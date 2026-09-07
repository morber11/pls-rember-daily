//go:build linux

#include "gtk_linux.h"

#include <gtk/gtk.h>
#include <gdk/x11/gdkx.h>

#if !GTK_CHECK_VERSION(4, 8, 0)
#error "Pls Rember Daily requires GTK 4.8 or newer"
#endif

static const int minimum_player_width = 480;
static const int minimum_player_height = 270;
static const int display_size_divisor = 4;

typedef struct
{
    int width;
    int height;
} player_size;

typedef struct
{
    GtkApplication *application;
    const char *video_path;
    GtkMediaStream *media;
    GtkWindow *window;
    char *error_message;
} playback_state;

static player_size calculate_player_size(int display_width, int display_height)
{
    player_size size = {
        .width = display_width / display_size_divisor,
        .height = display_height / display_size_divisor,
    };

    if (size.width < minimum_player_width ||
        size.height < minimum_player_height)
    {
        return (player_size){minimum_player_width, minimum_player_height};
    }

    if (size.width * minimum_player_height >
        size.height * minimum_player_width)
    {
        size.width = size.height * minimum_player_width /
                     minimum_player_height;
    }
    else
    {
        size.height = size.width * minimum_player_height /
                      minimum_player_width;
    }

    return size;
}

static void prepare_x11_player_window(GtkWidget *widget, gpointer user_data)
{
    (void)user_data;

    GdkDisplay *display = gtk_widget_get_display(widget);
    if (!GDK_IS_X11_DISPLAY(display))
    {
        return;
    }

    GdkMonitor *monitor = gdk_x11_display_get_primary_monitor(display);
    if (monitor == NULL)
    {
        return;
    }

    GdkRectangle geometry;
    gdk_monitor_get_geometry(monitor, &geometry);

    const player_size size = calculate_player_size(
        geometry.width, geometry.height);
    gtk_window_set_default_size(GTK_WINDOW(widget), size.width, size.height);

    GdkRectangle workarea;
    gdk_x11_monitor_get_workarea(monitor, &workarea);

    GdkSurface *surface = gtk_native_get_surface(GTK_NATIVE(widget));

    const int scale = gdk_surface_get_scale_factor(surface);
    const int x = (workarea.x + workarea.width - size.width) * scale;
    const int y = (workarea.y + workarea.height - size.height) * scale;

    Display *xdisplay = gdk_x11_display_get_xdisplay(display);
    Window xwindow = gdk_x11_surface_get_xid(surface);

    XSetWindowAttributes attributes = {
        .override_redirect = True,
    };

    XChangeWindowAttributes(xdisplay, xwindow, CWOverrideRedirect, &attributes);
    XMoveResizeWindow(xdisplay, xwindow, x, y,
                      size.width * scale, size.height * scale);
    XRaiseWindow(xdisplay, xwindow);
    XFlush(xdisplay);
}

static void size_player_window(GtkWidget *widget, gpointer user_data)
{
    (void)user_data;

    GdkDisplay *display = gtk_widget_get_display(widget);
    if (GDK_IS_X11_DISPLAY(display))
    {
        return;
    }

    GdkSurface *surface = gtk_native_get_surface(GTK_NATIVE(widget));
    GdkMonitor *monitor = gdk_display_get_monitor_at_surface(display, surface);

    if (monitor == NULL)
    {
        return;
    }

    GdkRectangle geometry;
    gdk_monitor_get_geometry(monitor, &geometry);

    const player_size size = calculate_player_size(
        geometry.width, geometry.height);
    gtk_window_set_default_size(GTK_WINDOW(widget), size.width, size.height);
}

static void close_player_window(playback_state *state)
{
    if (state->window != NULL)
    {
        GtkWindow *window = state->window;
        state->window = NULL;
        gtk_window_destroy(window);
    }
}

static void player_window_destroyed(GtkWidget *widget, gpointer user_data)
{
    (void)widget;

    playback_state *state = user_data;
    state->window = NULL;
}

static void media_ended(GObject *object, GParamSpec *parameter,
                        gpointer user_data)
{
    (void)parameter;
    playback_state *state = user_data;

    if (gtk_media_stream_get_ended(GTK_MEDIA_STREAM(object)))
    {
        close_player_window(state);
    }
}

static void media_failed(GObject *object, GParamSpec *parameter,
                         gpointer user_data)
{
    (void)parameter;
    playback_state *state = user_data;

    const GError *error = gtk_media_stream_get_error(GTK_MEDIA_STREAM(object));

    if (error == NULL)
    {
        return;
    }

    if (state->error_message == NULL)
    {
        state->error_message = g_strdup(error->message);
    }

    close_player_window(state);
}

static void activate_player(GtkApplication *application, gpointer user_data)
{
    playback_state *state = user_data;

    state->media = gtk_media_file_new_for_filename(state->video_path);

    GtkWidget *picture = gtk_picture_new_for_paintable(
        GDK_PAINTABLE(state->media));

    gtk_picture_set_content_fit(GTK_PICTURE(picture), GTK_CONTENT_FIT_CONTAIN);
    gtk_widget_set_hexpand(picture, TRUE);
    gtk_widget_set_vexpand(picture, TRUE);

    state->window = GTK_WINDOW(gtk_application_window_new(application));

    gtk_window_set_title(state->window, "Pls Rember Daily");
    gtk_window_set_default_size(
        state->window, minimum_player_width, minimum_player_height);

    gtk_window_set_decorated(state->window, FALSE);
    gtk_window_set_resizable(state->window, FALSE);
    gtk_window_set_child(state->window, picture);

    g_signal_connect(state->window, "realize",
                     G_CALLBACK(prepare_x11_player_window), NULL);
    g_signal_connect(state->window, "map",
                     G_CALLBACK(size_player_window), NULL);
    g_signal_connect(state->window, "destroy",
                     G_CALLBACK(player_window_destroyed), state);
    g_signal_connect(state->media, "notify::ended",
                     G_CALLBACK(media_ended), state);
    g_signal_connect(state->media, "notify::error",
                     G_CALLBACK(media_failed), state);

    const GError *error = gtk_media_stream_get_error(state->media);
    if (error != NULL)
    {
        state->error_message = g_strdup(error->message);
        close_player_window(state);
        return;
    }

    gtk_window_present(state->window);
    gtk_media_stream_play(state->media);
}

void *linux_application_new(const char *application_id)
{
#if GLIB_CHECK_VERSION(2, 74, 0)
    return gtk_application_new(application_id, G_APPLICATION_DEFAULT_FLAGS);
#else
    return gtk_application_new(application_id, G_APPLICATION_FLAGS_NONE);
#endif
}

int linux_application_register(void *application, char **error_message)
{
    GError *error = NULL;
    
    const gboolean registered = g_application_register(
        G_APPLICATION(application), NULL, &error);

    if (!registered)
    {
        *error_message = g_strdup(error->message);
        g_error_free(error);
    }

    return registered;
}

int linux_application_is_remote(void *application)
{
    return g_application_get_is_remote(G_APPLICATION(application));
}

void linux_application_release(void *application)
{
    g_object_unref(application);
}

char *linux_play_video(void *application, const char *video_path)
{
    playback_state state = {
        .application = GTK_APPLICATION(application),
        .video_path = video_path,
    };

    g_signal_connect(state.application, "activate",
                     G_CALLBACK(activate_player), &state);

    const int status = g_application_run(
        G_APPLICATION(state.application), 0, NULL);

    g_signal_handlers_disconnect_by_data(state.application, &state);

    if (state.media != NULL)
    {
        g_object_unref(state.media);
    }

    if (status != 0 && state.error_message == NULL)
    {
        state.error_message = g_strdup_printf(
            "GTK application exited with status %d", status);
    }

    return state.error_message;
}

void linux_error_free(char *error_message)
{
    g_free(error_message);
}
