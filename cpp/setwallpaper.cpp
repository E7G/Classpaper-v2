#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
#include <tchar.h>

BOOL CALLBACK EnumWindowsProc(HWND hwnd, LPARAM lParam);

//声明一个结构体用于存储获取到的所有窗口类名
typedef struct windows_class {
    char window_class_name[256];
    HWND win_hwnd;
    windows_class* next;
} windows_class;

//声明一个全局结构体
windows_class* class_name;

//记录屏幕窗口类数量
int num;

int main(int argc ,char **argv)
{
	if(argc!=2){
		printf("用法： 本程序 窗口程序名称 \n 用于将窗口程序嵌入桌面");
		return 0;
	}
    //获取窗口句柄
    HWND hWnd = FindWindow(_T("Progman"), _T("Program Manager"));
    if (hWnd == NULL) {
        printf("无法获取桌面句柄");
        //getchar();
        return 0;
    }

    //发送多屏消息
    SendMessage(hWnd, 0x052c, 0, 0);
    //SendMessage(hWnd, 0x052c, 0, 1);

    //结构体初始化
    class_name = (windows_class*)malloc(sizeof(windows_class));

    //枚举屏幕上所有窗口
    EnumWindows(EnumWindowsProc, 0);

    //循环比对找到WorkerW类
    for (int i = 0; i < num; ++i) {
        if (strncmp(class_name->window_class_name, "WorkerW", strlen(class_name->window_class_name)) == 0) {
            //以有效字符比对，防止连同字符“0”等无效字符也一同包含在一起比对
            HWND window_hwnd = FindWindowExA(class_name->win_hwnd, 0, "SHELLDLL_DefView", NULL);
            if (window_hwnd == NULL) {
                //无法获取句柄代表该workerw类窗口没有子窗口也就是获取到图标下面的WorkerW类窗口了
                //直接关闭该窗口
                SendMessage(class_name->win_hwnd, WM_CLOSE, 0, 0);
                break;
            }
            else {
                //获取成功看一下下一个窗口是不是Progman
                class_name = class_name->next;
                if (strcmp(class_name->window_class_name, "Progman") == 0) {
                    HWND window_hwnd = FindWindowExA(class_name->win_hwnd, 0, "WorkerW", NULL);
                    //获取图标下面的WorkerW子窗口
                    if (window_hwnd == NULL) {
                        //获取不到代表该窗口已经被关闭了
                        printf("该窗口已经被关闭..");
                        //getchar();
                        break;
                    }
                    else {
                        //结束窗口
                        SendMessage(window_hwnd, WM_CLOSE, 0, 0);
                    }
                }
                else {
                    //如果不是Progman就代表WorkerW类窗口的屏幕Z序列高于Progman
                    //就说明获取到了WorkerW类窗口，直接关闭即可
                    SendMessage(class_name->win_hwnd, WM_CLOSE, 0, 0);
                }
            }
        }
        class_name = class_name->next;
    }

    //到了这一步就可以将你的视频窗口嵌入到Progman类窗口当中了
    //嵌入之后也不会被遮挡，因为遮挡的WorkerW窗口已经被关闭了

    //所以这里我就随便嵌入一个窗口进入桌面看一下会不会被遮挡
    //获取要嵌入窗口的句柄
    HWND test_hwnd = FindWindow(NULL, _T(argv[1]));



    if (test_hwnd == NULL) {
        printf("要嵌入的窗口不存在..");
        //getchar();
        return -1;
    }
    else {
    
        SetParent(test_hwnd, hWnd);


        printf("窗口嵌入完成");
    }
        
	

    //getchar();
    return 0;
}

BOOL CALLBACK EnumWindowsProc(HWND hwnd, LPARAM lParam)
{
    //声明结构体
    windows_class* enum_calss_name;
    //初始化
    enum_calss_name = (windows_class*)malloc(sizeof(windows_class));
    //填充0到类名变量中
    memset(enum_calss_name->window_class_name, 0, sizeof(enum_calss_name->window_class_name));
    //获取窗口类名
    GetClassNameA(hwnd, enum_calss_name->window_class_name, sizeof(enum_calss_name->window_class_name));
    //获取窗口句柄
    enum_calss_name->win_hwnd = hwnd;
    //递增类数量
    num += 1;
    //链表形式存储
    enum_calss_name->next = class_name;
    class_name = enum_calss_name;
    return TRUE; //这里必须返回TRUE,返回FALSE就不在枚举了
}
