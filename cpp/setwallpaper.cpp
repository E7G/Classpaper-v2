#include <iostream>
#include <string>
#include <windows.h>
#include <tchar.h>
HWND workerw=NULL;     //第二个WorkerW窗口句柄
 
inline BOOL CALLBACK EnumWindowsProc1(HWND handle,LPARAM lparam)
{
    //获取第一个WorkerW窗口
    HWND defview = FindWindowEx(handle, 0, _T("SHELLDLL_DefView"), NULL);
 
    if (defview != NULL) //找到第一个WorkerW窗口
    {
        //获取第二个WorkerW窗口的窗口句柄
        workerw = FindWindowEx(0, handle, _T("WorkerW"), 0);
    }
    return true;
}
 
//参数myAppHwnd为你开发的窗口程序的窗口句柄
void SetDesktop(HWND myAppHwnd)
{
    int result;
    HWND windowHandle = FindWindow(_T("Progman"), NULL);
    SendMessageTimeout(windowHandle, 0x052c, 0 ,0, SMTO_NORMAL, 0x3e8,(PDWORD_PTR)&result);
 
    //枚举窗口
    EnumWindows(EnumWindowsProc1,(LPARAM)NULL);
 
    //隐藏第二个WorkerW窗口，当以Progman为父窗口时需要对其进行隐藏，
    //不然程序窗口会被第二个WorkerW覆盖
    ShowWindow(workerw,SW_HIDE);
 
    SetParent(myAppHwnd,windowHandle);
}


struct EnumWindowsData {
    std::string keyword;
    HWND hwnd = nullptr;
};

BOOL CALLBACK EnumWindowsProc2(HWND hwnd, LPARAM lParam) {
    EnumWindowsData* data = reinterpret_cast<EnumWindowsData*>(lParam);

    char buffer[256];
    GetWindowTextA(hwnd, buffer, sizeof(buffer));

    std::string title(buffer);

    if (title.find(data->keyword) != std::string::npos) {
        data->hwnd = hwnd;
        return FALSE;  // 停止枚举
    }

    return TRUE;
}

int main(int argc ,char **argv)
{
	if(argc!=2){
		printf("usage: program.exe window_name \n use to stick window into desktop.");
		return 0;
	}

    EnumWindowsData data;
    data.keyword = argv[1];

    EnumWindows(EnumWindowsProc2, reinterpret_cast<LPARAM>(&data));

    if (data.hwnd != nullptr) {
        std::cout << "found :  " << data.hwnd << std::endl;

        SetDesktop(data.hwnd);


        printf("success");
    } else {
        printf("failed ! not found..");
        //getchar();
        return -1;
    }

        
	

    //getchar();
    return 0;
}
