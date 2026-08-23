// go:build windows

#include "player_windows.h"

#include <windows.h>
#include <mfapi.h>
#include <mferror.h>
#include <mfmediaengine.h>

#include <new>

// uses Media Foundation COM object to play video in a top level window
namespace
{

    constexpr wchar_t player_window_class[] = L"PlsRemberDailyPlayer";
    constexpr wchar_t player_window_title[] = L"PlsRemberDaily";

    constexpr UINT playback_started_message = WM_APP + 1; // playback started
    constexpr UINT playback_ended_message = WM_APP + 2;   // playback ended
    constexpr UINT playback_error_message = WM_APP + 3;   // playback error

    constexpr int icon_resource_id = 1;
    constexpr int minimum_player_width = 480;
    constexpr int minimum_player_height = 270;
    constexpr int display_size_divisor = 4;

    struct player_size
    {
        int width;
        int height;
    };

    struct window_state
    {
        HRESULT result = S_OK;
    };

    player_size calculate_player_size(int display_width, int display_height)
    {
        player_size size{
            display_width / display_size_divisor,
            display_height / display_size_divisor};

        if (size.width < minimum_player_width ||
            size.height < minimum_player_height)
        {
            return {minimum_player_width, minimum_player_height};
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

    class player_callback final : public IMFMediaEngineNotify
    {
    public:
        explicit player_callback(HWND window) : window_(window) {}

        HRESULT STDMETHODCALLTYPE QueryInterface(REFIID interface_id,
                                                 void **object) override
        {
            if (!object)
            {
                return E_POINTER;
            }

            // more COM interface boiler plate
            // since the interface can also support IUnknown
            if (interface_id == IID_IUnknown ||
                interface_id == __uuidof(IMFMediaEngineNotify))
            {
                *object = static_cast<IMFMediaEngineNotify *>(this);
                AddRef(); // now take COM ownership, increment references. defined below
                return S_OK;
            }

            *object = nullptr;
            return E_NOINTERFACE;
        }

        // COM references use manual reference counting so we need to actually increment it ourself
        ULONG STDMETHODCALLTYPE AddRef() override
        {
            return static_cast<ULONG>(InterlockedIncrement(&reference_count_));
        }

        // Release / decrement when finished using COM object so it can be destroyed automatically
        ULONG STDMETHODCALLTYPE Release() override
        {
            const auto references =
                static_cast<ULONG>(InterlockedDecrement(&reference_count_));

            if (references == 0)
            {
                delete this;
            }

            return references;
        }

        HRESULT STDMETHODCALLTYPE EventNotify(DWORD event, DWORD_PTR,
                                              DWORD param2) override
        {
            if (event == MF_MEDIA_ENGINE_EVENT_ERROR)
            {
                PostMessageW(window_, playback_error_message, 0,
                             static_cast<LPARAM>(static_cast<HRESULT>(param2)));
                return S_OK;
            }

            if (event == MF_MEDIA_ENGINE_EVENT_PLAYING)
            {
                PostMessageW(window_, playback_started_message, 0, 0);
            }
            else if (event == MF_MEDIA_ENGINE_EVENT_ENDED)
            {
                PostMessageW(window_, playback_ended_message, 0, 0);
            }

            return S_OK;
        }

    private:
        ~player_callback() = default;

        LONG reference_count_ = 1;
        HWND window_ = nullptr;
    };

    // handles our message queue that we posted messages to earlier
    LRESULT CALLBACK player_window_proc(HWND window, UINT message, WPARAM w_param,
                                        LPARAM l_param)
    {
        auto *state = reinterpret_cast<window_state *>(
            GetWindowLongPtrW(window, GWLP_USERDATA));

        if (message == WM_NCCREATE)
        {
            const auto *creation = reinterpret_cast<CREATESTRUCTW *>(l_param);
            state = static_cast<window_state *>(creation->lpCreateParams);

            SetWindowLongPtrW(window, GWLP_USERDATA,
                              reinterpret_cast<LONG_PTR>(state));
        }

        switch (message)
        {
        case playback_started_message:
            ShowWindow(window, SW_SHOW);
            SetForegroundWindow(window);
            return 0;

        case playback_ended_message:
            DestroyWindow(window);
            return 0;

        case playback_error_message:
            if (state)
            {
                state->result = static_cast<HRESULT>(l_param);
            }
            DestroyWindow(window);
            return 0;

        case WM_SETCURSOR: // explicit cursor handling so we don't get loading/spinner cursor
            SetCursor(LoadCursor(nullptr, IDC_ARROW));
            return TRUE;

        case WM_CLOSE:
            DestroyWindow(window);
            return 0;

        case WM_DESTROY:
            PostQuitMessage(0);
            return 0;

        default:
            return DefWindowProcW(window, message, w_param, l_param);
        }
    }

    HRESULT register_player_window_class(HINSTANCE instance)
    {
        const auto icon = static_cast<HICON>(LoadImageW(
            instance, MAKEINTRESOURCEW(icon_resource_id), IMAGE_ICON, 0, 0,
            LR_DEFAULTSIZE | LR_SHARED));

        WNDCLASSEXW window_class{};

        window_class.cbSize = sizeof(window_class);
        window_class.hInstance = instance;
        window_class.lpszClassName = player_window_class;
        window_class.lpfnWndProc = player_window_proc;

        window_class.hIcon = icon;
        window_class.hIconSm = icon;

        window_class.hbrBackground = static_cast<HBRUSH>(GetStockObject(BLACK_BRUSH));

        if (RegisterClassExW(&window_class) ||
            GetLastError() == ERROR_CLASS_ALREADY_EXISTS)
        {
            return S_OK;
        }

        return HRESULT_FROM_WIN32(GetLastError());
    }

    HRESULT create_player_window(HINSTANCE instance, window_state *state,
                                 HWND *window)
    {
        RECT work_area{};
        if (!SystemParametersInfoW(SPI_GETWORKAREA, 0, &work_area, 0))
        {
            return HRESULT_FROM_WIN32(GetLastError());
        }

        const player_size size = calculate_player_size(
            GetSystemMetrics(SM_CXSCREEN),
            GetSystemMetrics(SM_CYSCREEN));

        const int x = work_area.right - size.width;
        const int y = work_area.bottom - size.height;

        *window = CreateWindowExW(
            WS_EX_APPWINDOW | WS_EX_TOPMOST,
            player_window_class,
            player_window_title,
            WS_POPUP,
            x,
            y,
            size.width,
            size.height,
            nullptr,
            nullptr,
            instance,
            state);

        if (!*window)
        {
            return HRESULT_FROM_WIN32(GetLastError());
        }

        return S_OK;
    }

} // namespace

extern "C" long play_video(const wchar_t *video_url)
{
    HRESULT result = CoInitializeEx(nullptr, COINIT_APARTMENTTHREADED);
    const bool should_uninitialize_com = SUCCEEDED(result);

    if (FAILED(result) && result != RPC_E_CHANGED_MODE) // if thread is already initialized by another apartment, continue
    {
        return result;
    }

    result = MFStartup(MF_VERSION);
    if (FAILED(result))
    {
        if (should_uninitialize_com)
        {
            CoUninitialize();
        }
        return result;
    }

    const auto instance = GetModuleHandleW(nullptr);
    result = register_player_window_class(instance);

    window_state state{};
    HWND window = nullptr;
    if (SUCCEEDED(result))
    {
        result = create_player_window(instance, &state, &window);
    }

    auto *callback =
        SUCCEEDED(result)
            ? new (std::nothrow) player_callback(window)
            : nullptr;

    if (SUCCEEDED(result) && !callback)
    {
        result = E_OUTOFMEMORY;
    }

    IMFAttributes *attributes = nullptr;
    // checking SUCCEEDED(result) is ridiculoussly ugly
    // but it works and is straightforward enough to follow
    if (SUCCEEDED(result))
    {
        result = MFCreateAttributes(&attributes, 2);
    }

    if (SUCCEEDED(result))
    {
        // use our callback for notifications
        result = attributes->SetUnknown(MF_MEDIA_ENGINE_CALLBACK, callback);
    }

    if (SUCCEEDED(result))
    {
        result = attributes->SetUINT64(
            MF_MEDIA_ENGINE_PLAYBACK_HWND,
            // reinterpret_cast is needed because HWND
            // is pointer like handle but api wants 64 bit int
            reinterpret_cast<UINT64>(window));
    }

    IMFMediaEngineClassFactory *factory = nullptr;
    if (SUCCEEDED(result))
    {
        result = CoCreateInstance(
            CLSID_MFMediaEngineClassFactory,
            nullptr,
            CLSCTX_INPROC_SERVER,
            IID_PPV_ARGS(&factory));
    }

    IMFMediaEngine *player = nullptr;
    if (SUCCEEDED(result))
    {
        result = factory->CreateInstance(0, attributes, &player);
    }

    // magic COM string type similar to wchar_t
    BSTR source = nullptr;
    if (SUCCEEDED(result))
    {
        source = SysAllocString(video_url);
        if (!source)
        {
            result = E_OUTOFMEMORY;
        }
    }

    if (SUCCEEDED(result))
    {
        result = player->SetSource(source);
    }

    if (SUCCEEDED(result))
    {
        result = player->Play();
    }

    if (SUCCEEDED(result))
    {
        MSG message{};
        while (GetMessageW(&message, nullptr, 0, 0) > 0)
        {
            TranslateMessage(&message);
            DispatchMessageW(&message);
        }
        result = state.result;
    }

    if (player)
    {
        player->Shutdown();
        player->Release();
    }

    if (source)
    {
        SysFreeString(source);
    }

    if (factory)
    {
        factory->Release();
    }

    if (attributes)
    {
        attributes->Release();
    }

    if (callback)
    {
        callback->Release();
    }

    if (window && IsWindow(window))
    {
        DestroyWindow(window);
    }

    MFShutdown();
    if (should_uninitialize_com)
    {
        CoUninitialize();
    }

    return result;
}
