#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
#include <tchar.h>

BOOL CALLBACK EnumWindowsProc(HWND hwnd, LPARAM lParam);

//����һ���ṹ�����ڴ洢��ȡ�������д�������
typedef struct windows_class {
    char window_class_name[256];
    HWND win_hwnd;
    windows_class* next;
} windows_class;

//����һ��ȫ�ֽṹ��
windows_class* class_name;

//��¼��Ļ����������
int num;

int main(int argc ,char **argv)
{
	if(argc!=2){
		printf("�÷��� ������ ���ڳ������� \n ���ڽ����ڳ���Ƕ������");
		return 0;
	}
    //��ȡ���ھ��
    HWND hWnd = FindWindow(_T("Progman"), _T("Program Manager"));
    if (hWnd == NULL) {
        printf("�޷���ȡ������");
        //getchar();
        return 0;
    }

    //���Ͷ�����Ϣ
    SendMessage(hWnd, 0x052c, 0, 0);
    //SendMessage(hWnd, 0x052c, 0, 1);

    //�ṹ���ʼ��
    class_name = (windows_class*)malloc(sizeof(windows_class));

    //ö����Ļ�����д���
    EnumWindows(EnumWindowsProc, 0);

    //ѭ���ȶ��ҵ�WorkerW��
    for (int i = 0; i < num; ++i) {
        if (strncmp(class_name->window_class_name, "WorkerW", strlen(class_name->window_class_name)) == 0) {
            //����Ч�ַ��ȶԣ���ֹ��ͬ�ַ���0������Ч�ַ�Ҳһͬ������һ��ȶ�
            HWND window_hwnd = FindWindowExA(class_name->win_hwnd, 0, "SHELLDLL_DefView", NULL);
            if (window_hwnd == NULL) {
                //�޷���ȡ��������workerw�ര��û���Ӵ���Ҳ���ǻ�ȡ��ͼ�������WorkerW�ര����
                //ֱ�ӹرոô���
                SendMessage(class_name->win_hwnd, WM_CLOSE, 0, 0);
                break;
            }
            else {
                //��ȡ�ɹ���һ����һ�������ǲ���Progman
                class_name = class_name->next;
                if (strcmp(class_name->window_class_name, "Progman") == 0) {
                    HWND window_hwnd = FindWindowExA(class_name->win_hwnd, 0, "WorkerW", NULL);
                    //��ȡͼ�������WorkerW�Ӵ���
                    if (window_hwnd == NULL) {
                        //��ȡ��������ô����Ѿ����ر���
                        printf("�ô����Ѿ����ر�..");
                        //getchar();
                        break;
                    }
                    else {
                        //��������
                        SendMessage(window_hwnd, WM_CLOSE, 0, 0);
                    }
                }
                else {
                    //�������Progman�ʹ���WorkerW�ര�ڵ���ĻZ���и���Progman
                    //��˵����ȡ����WorkerW�ര�ڣ�ֱ�ӹرռ���
                    SendMessage(class_name->win_hwnd, WM_CLOSE, 0, 0);
                }
            }
        }
        class_name = class_name->next;
    }

    //������һ���Ϳ��Խ������Ƶ����Ƕ�뵽Progman�ര�ڵ�����
    //Ƕ��֮��Ҳ���ᱻ�ڵ�����Ϊ�ڵ���WorkerW�����Ѿ����ر���

    //���������Ҿ����Ƕ��һ�����ڽ������濴һ�»᲻�ᱻ�ڵ�
    //��ȡҪǶ�봰�ڵľ��
    HWND test_hwnd = FindWindow(NULL, _T(argv[1]));



    if (test_hwnd == NULL) {
        printf("ҪǶ��Ĵ��ڲ�����..");
        //getchar();
        return -1;
    }
    else {
    
        SetParent(test_hwnd, hWnd);


        printf("����Ƕ�����");
    }
        
	

    //getchar();
    return 0;
}

BOOL CALLBACK EnumWindowsProc(HWND hwnd, LPARAM lParam)
{
    //�����ṹ��
    windows_class* enum_calss_name;
    //��ʼ��
    enum_calss_name = (windows_class*)malloc(sizeof(windows_class));
    //���0������������
    memset(enum_calss_name->window_class_name, 0, sizeof(enum_calss_name->window_class_name));
    //��ȡ��������
    GetClassNameA(hwnd, enum_calss_name->window_class_name, sizeof(enum_calss_name->window_class_name));
    //��ȡ���ھ��
    enum_calss_name->win_hwnd = hwnd;
    //����������
    num += 1;
    //������ʽ�洢
    enum_calss_name->next = class_name;
    class_name = enum_calss_name;
    return TRUE; //������뷵��TRUE,����FALSE�Ͳ���ö����
}
