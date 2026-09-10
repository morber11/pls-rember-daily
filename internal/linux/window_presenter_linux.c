//go:build linux

#include "window_presenter_linux.h"

static const int display_size_divisor = 4;

static gboolean reposition_x11_player_window(gpointer user_data)
{
    GtkWidget *widget = GTK_WIDGET(user_data);
    prepare_x11_player_window(widget);
    g_object_unref(widget);
    return G_SOURCE_REMOVE;
}

player_size calculate_player_size(int display_width, int display_height)
{
    player_size size = {
        .width = display_width / display_size_divisor,
        .height = display_height / display_size_divisor,
    };

    if (size.width < player_minimum_width ||
        size.height < player_minimum_height)
    {
        return (player_size){player_minimum_width, player_minimum_height};
    }

    if (size.width * player_minimum_height >
        size.height * player_minimum_width)
    {
        size.width = size.height * player_minimum_width /
                     player_minimum_height;
    }
    else
    {
        size.height = size.width * player_minimum_height /
                      player_minimum_width;
    }

    return size;
}

void prepare_wayland_player_window(GtkWidget *widget)
{
    GdkDisplay *display = gtk_widget_get_display(widget);
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

void prepare_player_window(GtkWidget *widget, gpointer user_data)
{
    (void)user_data;
    prepare_x11_player_window(widget);
}

void size_player_window(GtkWidget *widget, gpointer user_data)
{
    (void)user_data;

    if (prepare_x11_player_window(widget))
    {
        g_idle_add(reposition_x11_player_window, g_object_ref(widget));
        return;
    }

    prepare_wayland_player_window(widget);
}
