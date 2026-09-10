//go:build linux

#include "gtk_linux.h"
#include "window_presenter_linux.h"

#include <gtk/gtk.h>

#if !GTK_CHECK_VERSION(4, 8, 0)
#error "Pls Rember Daily requires GTK 4.8 or newer"
#endif

typedef struct
{
    GtkApplication *application;
    const char *video_path;
    GtkMediaStream *media;
    GtkWindow *window;
    char *error_message;
} playback_state;

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
        state->window, player_minimum_width, player_minimum_height);

    gtk_window_set_decorated(state->window, FALSE);
    gtk_window_set_resizable(state->window, FALSE);
    gtk_window_set_child(state->window, picture);

    g_signal_connect(state->window, "realize",
                     G_CALLBACK(prepare_player_window), state);
    g_signal_connect(state->window, "map",
                     G_CALLBACK(size_player_window), state);
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
    configure_linux_rendering();

    if (g_getenv("GSK_RENDERER") == NULL)
    {
        g_setenv("GSK_RENDERER", "cairo", FALSE);
    }

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

    if (state.application != NULL)
    {
        g_signal_handlers_disconnect_by_data(state.application, &state);
    }

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
