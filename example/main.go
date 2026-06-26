package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	pritunl "github.com/zhong/pritunl-go-sdk"
)

func main() {
	baseURL := flag.String("base", "https://localhost", "Pritunl base URL")
	token := flag.String("token", os.Getenv("PRITUNL_API_TOKEN"), "API token")
	secret := flag.String("secret", os.Getenv("PRITUNL_API_SECRET"), "API secret")
	insecure := flag.Bool("insecure", true, "skip TLS verification")
	cleanup := flag.Bool("cleanup", false, "delete test server, user and organization after creation")
	serverName := flag.String("server", "go-sdk-test-server", "test server name")
	orgName := flag.String("org", "go-sdk-test-org", "test organization name")
	userName := flag.String("user", "go-sdk-test-user", "test user name")
	flag.Parse()

	if *token == "" || *secret == "" {
		log.Fatal("请提供 -token 和 -secret，或设置环境变量 PRITUNL_API_TOKEN / PRITUNL_API_SECRET")
	}

	client := pritunl.NewClient(*baseURL, *token, *secret, *insecure)

	// 1. 获取服务器状态
	status, err := client.GetStatus()
	if err != nil {
		log.Fatalf("获取状态失败: %v", err)
	}
	fmt.Printf("Pritunl 状态: version=%s, orgs=%d, users=%d, servers=%d, online=%d\n",
		status.ServerVersion, status.OrgCount, status.UserCount, status.ServerCount, status.UsersOnline)

	// 2. 创建启用 OTP（多因子鉴别）的 VPN Server
	server, err := client.CreateServer(pritunl.CreateServerRequest{
		Name:     *serverName,
		Protocol: "udp",
		Cipher:   "aes256",
		Hash:     "sha256",
		OTPAuth:  true,
	})
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}
	fmt.Printf("创建服务器成功: %s (%s), otp_auth=%v\n", server.Name, server.ID, server.OTPAuth)

	// 3. 创建测试组织
	org, err := client.CreateOrganization(*orgName)
	if err != nil {
		log.Fatalf("创建组织失败: %v", err)
	}
	fmt.Printf("创建组织成功: %s (%s)\n", org.Name, org.ID)

	// 4. 把组织关联到服务器
	if err := client.AttachOrganization(server.ID, org.ID); err != nil {
		log.Fatalf("关联服务器和组织失败: %v", err)
	}
	fmt.Println("组织已关联到服务器")

	// 5. 在组织中创建测试用户
	user, err := client.CreateUser(org.ID, pritunl.CreateUserRequest{
		Name:  *userName,
		Email: "test@example.com",
	})
	if err != nil {
		log.Fatalf("创建用户失败: %v", err)
	}
	fmt.Printf("创建用户成功: %s (%s)\n", user.Name, user.ID)

	// 6. 生成用户 OTP Secret（多因子鉴别绑定用）
	user, err = client.GenerateUserOTPSecret(org.ID, user.ID)
	if err != nil {
		log.Fatalf("生成 OTP Secret 失败: %v", err)
	}
	fmt.Printf("用户 OTP Secret 已生成: %s\n", user.OTPSecret)

	// 7. 获取用户配置下载链接和 MFA 绑定地址
	link, err := client.GetKeyLink(org.ID, user.ID)
	if err != nil {
		log.Fatalf("获取密钥链接失败: %v", err)
	}

	fmt.Println("\n=== 用户配置与多因子鉴别信息 ===")
	fmt.Printf("客户端导入 URI: %s\n", link.FullURI(*baseURL))
	fmt.Printf("配置下载地址:\n")
	fmt.Printf("  TAR: %s\n", joinURL(*baseURL, link.KeyURL))
	fmt.Printf("  ZIP: %s\n", joinURL(*baseURL, link.KeyZipURL))
	fmt.Printf("  ONC: %s\n", joinURL(*baseURL, link.KeyOncURL))
	fmt.Printf("Web 绑定/查看地址（含 MFA QR Code）: %s\n", joinURL(*baseURL, link.ViewURL))
	fmt.Printf("OTP Secret (用于手动绑定): %s\n", user.OTPSecret)
	fmt.Println("==================================")

	// 8. 下载 zip 配置归档
	archive, err := client.DownloadKeyArchive(org.ID, user.ID, "zip")
	if err != nil {
		log.Fatalf("下载密钥归档失败: %v", err)
	}
	filename := fmt.Sprintf("%s_%s.zip", org.Name, user.Name)
	if err := pritunl.SaveUserKeyArchive(archive, filename); err != nil {
		log.Fatalf("保存密钥归档失败: %v", err)
	}
	fmt.Printf("密钥归档已保存到: %s (%d bytes)\n", filename, len(archive))
	if len(archive) < 100 {
		fmt.Println("  警告：归档文件过小，可能服务器尚未启动或用户未分配虚拟 IP")
	}

	// 9. 可选清理
	if *cleanup {
		if err := client.DeleteUser(org.ID, user.ID); err != nil {
			log.Printf("删除用户失败: %v", err)
		} else {
			fmt.Println("测试用户已删除")
		}
		if err := client.DeleteOrganization(org.ID); err != nil {
			log.Printf("删除组织失败: %v", err)
		} else {
			fmt.Println("测试组织已删除")
		}
		if err := client.DeleteServer(server.ID); err != nil {
			log.Printf("删除服务器失败: %v", err)
		} else {
			fmt.Println("测试服务器已删除")
		}
	} else {
		fmt.Println("已保留测试服务器、组织和用户，可登录 Web 控制台验证")
		fmt.Printf("  服务器: %s (%s)\n", server.Name, server.ID)
		fmt.Printf("  组织:   %s (%s)\n", org.Name, org.ID)
		fmt.Printf("  用户:   %s (%s)\n", user.Name, user.ID)
	}
}

// joinURL combines a base URL with a path. Very small helper for display only.
func joinURL(base, path string) string {
	if len(base) > 0 && base[len(base)-1] == '/' {
		base = base[:len(base)-1]
	}
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return base + path
}
