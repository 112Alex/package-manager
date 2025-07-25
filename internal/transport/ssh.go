package transport

import (
    "fmt"
    "io"
    "net"
    "os"
    "path"
    "time"

    "golang.org/x/crypto/ssh"
    "github.com/pkg/sftp"
)

type SSHConfig struct {
    Host     string // host:port
    User     string
    Password string // optional
    KeyPath  string // optional path to private key
}

type Client struct {
    sftp *sftp.Client
}

// Connect establishes SSH + SFTP session.
func Connect(cfg SSHConfig) (*Client, error) {
    auths := []ssh.AuthMethod{}
    if cfg.KeyPath != "" {
        key, err := os.ReadFile(cfg.KeyPath)
        if err != nil {
            return nil, err
        }
        signer, err := ssh.ParsePrivateKey(key)
        if err != nil {
            return nil, err
        }
        auths = append(auths, ssh.PublicKeys(signer))
    } else if cfg.Password != "" {
        auths = append(auths, ssh.Password(cfg.Password))
    } else {
        return nil, fmt.Errorf("no auth method provided")
    }

    if _, _, err := net.SplitHostPort(cfg.Host); err != nil {
        cfg.Host = net.JoinHostPort(cfg.Host, "22")
    }

    sshConf := &ssh.ClientConfig{
        User:            cfg.User,
        Auth:            auths,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // for test task; DO NOT use in prod
        Timeout:         10 * time.Second,
    }
    conn, err := ssh.Dial("tcp", cfg.Host, sshConf)
    if err != nil {
        return nil, err
    }
    sftpClient, err := sftp.NewClient(conn)
    if err != nil {
        return nil, err
    }
    return &Client{sftp: sftpClient}, nil
}

func (c *Client) Close() error {
    return c.sftp.Close()
}

// Upload uploads localPath to remoteDir (creates dir) keeping filename.
func (c *Client) Upload(localPath, remoteDir string) error {
    src, err := os.Open(localPath)
    if err != nil {
        return err
    }
    defer src.Close()

    base := path.Base(localPath)
    if err := c.sftp.MkdirAll(remoteDir); err != nil {
        return err
    }
    dst, err := c.sftp.Create(path.Join(remoteDir, base))
    if err != nil {
        return err
    }
    defer dst.Close()

    _, err = io.Copy(dst, src)
    return err
}

// Download downloads remotePath to localDir.
func (c *Client) Download(remotePath, localDir string) (string, error) {
    src, err := c.sftp.Open(remotePath)
    if err != nil {
        return "", err
    }
    defer src.Close()

    base := path.Base(remotePath)
    localPath := path.Join(localDir, base)
    dst, err := os.Create(localPath)
    if err != nil {
        return "", err
    }
    defer dst.Close()

    _, err = io.Copy(dst, src)
    return localPath, err
}
