//go:build linux

#include "window_presenter_linux.h"

#ifdef GDK_WINDOWING_X11
#include <gdk/x11/gdkx.h>

void configure_linux_rendering(void)
{
    if (g_getenv("GDK_DISABLE") != NULL)
    {
        return;
    }

    const char *backend = g_getenv("GDK_BACKEND");
    if (backend != NULL && g_strcmp0(backend, "x11") != 0)
    {
        return;
    }

    if (backend == NULL && g_getenv("WAYLAND_DISPLAY") != NULL)
    {
        return;
    }

    const char *display_name = g_getenv("DISPLAY");
    if (display_name == NULL)
    {
        return;
    }

    Display *xdisplay = XOpenDisplay(display_name);
    if (xdisplay == NULL)
    {
        return;
    }

    int major_opcode;
    int first_event;
    int first_error;

    if (!XQueryExtension(xdisplay, "DRI2", &major_opcode, &first_event,
                         &first_error))
    {
        g_setenv("GDK_DISABLE", "gl", FALSE);
    }

    XCloseDisplay(xdisplay);
}

gboolean prepare_x11_player_window(GtkWidget *widget)
{
    GdkDisplay *display = gtk_widget_get_display(widget);
    if (display == NULL || !GDK_IS_X11_DISPLAY(display))
    {
        return FALSE;
    }

    GdkMonitor *monitor = gdk_x11_display_get_primary_monitor(display);
    if (monitor == NULL)
    {
        return FALSE;
    }

    GdkRectangle geometry;
    gdk_monitor_get_geometry(monitor, &geometry);

    const player_size size = calculate_player_size(
        geometry.width, geometry.height);
    gtk_window_set_default_size(GTK_WINDOW(widget), size.width, size.height);

    GdkRectangle workarea;
    gdk_x11_monitor_get_workarea(monitor, &workarea);

    GdkSurface *surface = gtk_native_get_surface(GTK_NATIVE(widget));
    if (surface == NULL)
    {
        return FALSE;
    }

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

    return TRUE;
}
#else
void configure_linux_rendering(void)
{
}

gboolean prepare_x11_player_window(GtkWidget *widget)
{
    (void)widget;
    return FALSE;
}
#endif
