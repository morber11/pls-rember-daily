#ifndef PLS_REMBER_DAILY_GTK_LINUX_H
#define PLS_REMBER_DAILY_GTK_LINUX_H

void *linux_application_new(const char *application_id);
int linux_application_register(void *application, char **error_message);
int linux_application_is_remote(void *application);
void linux_application_release(void *application);

char *linux_play_video(void *application, const char *video_path);
void linux_error_free(char *error_message);

#endif
