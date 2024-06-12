#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
#include <tchar.h>
#include <Shobjidl.h>

extern "C" {const GUID CLSID_TaskbarList = {0x56FDF344, 0xFD6D, 0x11D0, 
{0x95, 0x8A, 0x00, 0x60, 0x97, 0xC9, 0xA0, 0x90} };
                          const GUID IID_ITaskbarList  = {0x56FDF342, 0xFD6D, 
0x11D0, {0x95, 0x8A, 0x00, 0x60, 0x97, 0xC9, 0xA0, 0x90} }; }

extern "C" {const GUID IID_ITaskbarList2 = {0x602D4995, 0xB13A, 0x429b, {0xA6, 
0x6E, 0x19, 0x35, 0xE4, 0x4F, 0x43, 0x17} }; }

extern "C" {const GUID IID_ITaskbarList3 = { 0xEA1AFB91, 0x9E28, 0x4B86, {0x90, 
0xE9, 0x9E, 0x9F, 0x8A, 0x5E, 0xEF, 0xAF} };}


void UnregisterTab(HWND tab)
{
    LPVOID lp = NULL ;
    CoInitialize(lp);//初始化COM库：没有这两句隐藏不起作用
    
    ITaskbarList3 *taskbar;

    if(!(tab))
        return;

    if(S_OK != CoCreateInstance(CLSID_TaskbarList, 0, CLSCTX_INPROC_SERVER, IID_ITaskbarList3, (void**)&taskbar))
        return;
    taskbar->HrInit();
    taskbar->UnregisterTab(tab);
    taskbar->Release();
}


BOOL ShowInTaskbar(HWND hWnd, BOOL bShow)
{
    LPVOID lp = NULL ;
    CoInitialize(lp);//初始化COM库：没有这两句隐藏不起作用

    HRESULT hr;
    ITaskbarList* pTaskbarList;
    hr = CoCreateInstance( CLSID_TaskbarList, NULL, CLSCTX_INPROC_SERVER, 
              IID_ITaskbarList, (void**)&pTaskbarList );

    if(SUCCEEDED(hr))
    {
         pTaskbarList->HrInit();
         if(bShow)
               pTaskbarList->AddTab(hWnd);
          else
               pTaskbarList->DeleteTab(hWnd);

         pTaskbarList->Release();
         return TRUE;
    }
 
 return FALSE;
}

void deltab(HWND test_hwnd){
       ShowWindow(test_hwnd, SW_HIDE);

       DWORD dwExStyle = GetWindowLong(test_hwnd, GWL_EXSTYLE);
     
       //dwExStyle &= ~(WS_VISIBLE);
     
       dwExStyle |= WS_EX_TOOLWINDOW;   

       dwExStyle &= ~(WS_EX_APPWINDOW);
     
       SetWindowLong(test_hwnd, GWL_EXSTYLE, dwExStyle);
     
       ShowWindow(test_hwnd, SW_SHOW);
     
     
       UpdateWindow(test_hwnd);


       ShowInTaskbar(test_hwnd, FALSE);

       UnregisterTab(test_hwnd);

 
      return;
}


int main(int argc ,char **argv)
{
	if(argc!=2){
		printf("用法： 本程序 窗口程序名称 \n 用于窗口程序");
		return 0;
	}
 
    //到了这一步就可以将你的视频窗口嵌入到Progman类窗口当中了
    //嵌入之后也不会被遮挡，因为遮挡的WorkerW窗口已经被关闭了

    //所以这里我就随便嵌入一个窗口进入桌面看一下会不会被遮挡
    //获取要嵌入窗口的句柄

    HWND test_hwnd = FindWindow(NULL, _T(argv[1]));


    if (test_hwnd == NULL) {
        printf("窗口不存在..");
        //getchar();
        //return -1;
    }
    else {
        deltab(test_hwnd);
        printf("完成");
        return 0;
    }
        
	
    test_hwnd = FindWindow(_T("Program Manager"), _T(argv[1]));


    if (test_hwnd == NULL) {
        printf("窗口不存在..");
        //getchar();
        return -1;
    }
    else {
        deltab(test_hwnd);
        printf("完成");
    }
        
    test_hwnd = FindWindow(_T(argv[1]), _T("Chrome Legacy Window"));


    if (test_hwnd == NULL) {
        printf("窗口不存在..");
        //getchar();
        return -1;
    }
    else {
        deltab(test_hwnd);
        printf("完成");
    }

    //getchar();
    return 0;
}

