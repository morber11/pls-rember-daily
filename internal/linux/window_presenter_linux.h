#ifndef PLS_REMBER_DAILY_WINDOW_PRESENTER_LINUX_H
#define PLS_REMBER_DAILY_WINDOW_PRESENTER_LINUX_H

#include <gtk/gtk.h>

enum
{
    player_minimum_width = 480,
    player_minimum_height = 270,
};

typedef struct
{
    int width;
    int height;
} player_size;

player_size calculate_player_size(int display_width, int display_height);
void configure_linux_rendering(void);
gboolean prepare_x11_player_window(GtkWidget *widget);
void prepare_wayland_player_window(GtkWidget *widget);
void prepare_player_window(GtkWidget *widget, gpointer user_data);
void size_player_window(GtkWidget *widget, gpointer user_data);

#endif
