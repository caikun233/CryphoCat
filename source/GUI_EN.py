import os
import tkinter as tk
from tkinter import messagebox, scrolledtext
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.backends import default_backend
import tkinter.filedialog as filedialog
import pyperclip

if not os.path.exists('RSAkeys'):
    print('Detected no `RSAkeys` folder, creating it automatically.\n')
    os.makedirs('RSAkeys')
if not os.path.exists('RSAkeys\\friend'):
    os.makedirs('RSAkeys\\friend')
if not os.path.exists('RSAkeys\\my'):
    os.makedirs('RSAkeys\\my')

if not os.path.exists('RSAkeys\\my\\public.txt'):
    open('RSAkeys\\my\\public.txt', 'w').close()
if not os.path.exists('RSAkeys\\my\\private.txt'):
    open('RSAkeys\\my\\private.txt', 'w').close()




def generate_key_pair():
    if not os.path.exists('RSAkeys'):
        os.makedirs('RSAkeys')
    if not os.path.exists('RSAkeys/my'):
        os.makedirs('RSAkeys/my')

    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048,
        backend=default_backend()
    )

    private_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption()
    )

    public_pem = private_key.public_key().public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )

    with open('RSAkeys/my/private.txt', 'wb') as f:
        f.write(private_pem)
    with open('RSAkeys/my/public.txt', 'wb') as f:
        f.write(public_pem)

    
    output_text.insert(tk.END, "RSA key pair generated successfully! Saved to RSAkeys/my/public.txt and RSAkeys/my/private.txt.\n\n")




def import_public_key():
    # 弹出文件选择窗口
    file_path = filedialog.askopenfilename(title="Select public key file.", filetypes=[("Text Files", "*.txt"), ("PEM Files", "*.pem")])
    if file_path:
        try:
            # 读取选择的公钥文件
            with open(file_path, 'rb') as f:
                public_pem = f.read()
                with open('RSAkeys/friend/public.txt', 'wb') as fw:
                    fw.write(public_pem)
            messagebox.showinfo("Prompt", "Public key imported successfully!")
            import_btn.config(text="Public Key Imported", state=tk.DISABLED)
            root.update()  # 更新窗口，立即生效
        except Exception as e:
            messagebox.showerror("Error", f"Public key import failed.：{str(e)}")




def encrypt_data():
    # 检查对方公钥文件是否存在
    if not os.path.isfile('RSAkeys/friend/public.txt'):
        response = messagebox.askquestion("Prompt", "Your friend's public key was not found. Do you want to use your own public key for encryption?")
        if response == 'yes':
            # 使用自己的公钥加密
            with open('RSAkeys/my/public.txt', 'rb') as f:
                public_pem = f.read()
                public_key = serialization.load_pem_public_key(public_pem, backend=default_backend())
        else:
            return  # 取消本次加密请求
    else:
        # 读取对方公钥
        with open('RSAkeys/friend/public.txt', 'rb') as f:
            public_pem = f.read()
            public_key = serialization.load_pem_public_key(public_pem, backend=default_backend())

    # 获取要加密的数据
    data = text_entry.get("1.0", tk.END).encode('utf-8')

    # 分段加密
    max_len = 127  # 限制每个分段的最大长度
    ciphertext_segments = []  # 定义一个空的列表用于存储加密后的分段
    if len(data) > max_len:
        segments = [data[i:i + max_len] for i in range(0, len(data), max_len)]
        for segment in segments:
            ciphertext_segment = public_key.encrypt(segment,
                                                   padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                                algorithm=hashes.SHA256(), label=None))
            ciphertext_segments.append(ciphertext_segment)
        ciphertext = b''.join(ciphertext_segments)
    else:
        ciphertext = public_key.encrypt(data,
                                        padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                     algorithm=hashes.SHA256(), label=None))
        ciphertext_segments.append(ciphertext)  # 添加单个密文到列表中

    # 将密文复制到剪贴板
    pyperclip.copy(ciphertext.hex())

    # 在输出文本框中显示密文
    output_text.delete("1.0", tk.END)  # 清空输出文本框
    output_text.insert(tk.END, f"Text after encrypt (copied automatically) ：{ciphertext.hex()}\n\n")
    output_text.see(tk.END)


def decrypt_data():
    # 读取自己的私钥
    with open('RSAkeys/my/private.txt', 'rb') as f:
        private_pem = f.read()
        private_key = serialization.load_pem_private_key(private_pem, password=None, backend=default_backend())

    # 获取要解密的数据
    ciphertext_hex = text_entry.get("1.0", tk.END).strip()  # 去除前后的空白字符

    try:
        ciphertext = bytes.fromhex(ciphertext_hex)
    except ValueError:
        messagebox.showerror("Error", "Invalid HEX ciphertext!")
        return

    # 分段解密
    if len(ciphertext) > 256:
        segments = [ciphertext[i: i + 256] for i in range(0, len(ciphertext), 256)]
        data_segments = []
        for segment in segments:
            data_segments.append(private_key.decrypt(segment,
                                                     padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                                  algorithm=hashes.SHA256(), label=None)))
        data = b"".join(data_segments)
    else:
        data = private_key.decrypt(ciphertext,
                                   padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                algorithm=hashes.SHA256(), label=None))

    output_text.insert(tk.END, f"Text after decrypt: {data.decode('utf-8')}\n\n")
    output_text.see(tk.END)


# 创建GUI窗口
root = tk.Tk()
root.geometry("800x600")  # 设置窗口大小
root.title("ChyphoCat")

# 生成RSA密钥对按钮
generate_btn = tk.Button(root, text="Generate RSA key pair", command=generate_key_pair)
generate_btn.place(x=100, y=30, width=200, height=50)  # 按钮位于 (100, 50) 的位置

# 导入公钥按钮
import_btn = tk.Button(root, text="Import Public Key", command=import_public_key)
import_btn.place(x=500, y=30, width=200, height=50)  # 按钮位于 (100, 50) 的位置

# 文本输入框
text_entry = scrolledtext.ScrolledText(root)
text_entry.place(x=80, y=100, width=640, height=190)  # 输入框位于 (100, 150) 的位置

# 加密按钮
encrypt_btn = tk.Button(root, text="Encrypt", command=encrypt_data)
encrypt_btn.place(x=100, y=320, width=200, height=50)

# 解密按钮
decrypt_btn = tk.Button(root, text="Decrypt", command=decrypt_data)
decrypt_btn.place(x=500, y=320, width=200, height=50)

# 输出文本框
output_text = scrolledtext.ScrolledText(root)
output_text.place(x=80, y=380, width=640, height=190)

# 隐藏控制台黑框
#root.wm_attributes('-topmost', 1)
#root.after(1, lambda: root.focus_force())

# 启动主循环
root.mainloop()
