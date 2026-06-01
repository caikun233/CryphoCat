import os
import tkinter as tk
from tkinter import messagebox, scrolledtext, simpledialog
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.backends import default_backend
import tkinter.filedialog as filedialog
import pyperclip

# -- 密钥存储路径 ---------------------------------------------------------------
_BASE   = 'RSAkeys'
_MY_DIR = os.path.join(_BASE, 'my')
_FR_DIR = os.path.join(_BASE, 'friend')
_MY_PRI = os.path.join(_MY_DIR, 'private.pem')
_MY_PUB = os.path.join(_MY_DIR, 'public.pem')
_FR_PUB = os.path.join(_FR_DIR, 'public.pem')

for _d in (_MY_DIR, _FR_DIR):
    os.makedirs(_d, exist_ok=True)

# -- 复古终端主题 ---------------------------------------------------------------
_BG   = '#0d0d0d'
_FG   = '#33cc33'
_OFGC = '#22aa22'
_BFBG = '#1c1c1c'
_BFAC = '#003300'
_FONT = ('Courier', 9)
_BFNT = ('Courier', 9, 'bold')
_SFNT = ('Courier', 8)


# -- 工具函数 -------------------------------------------------------------------

def _log(msg):
    output_text.configure(state=tk.NORMAL)
    output_text.insert(tk.END, msg + '\n')
    output_text.see(tk.END)
    output_text.configure(state=tk.DISABLED)


def _status(msg):
    status_var.set(msg)


def _ask_new_passphrase():
    """提示设置新口令（需二次确认）。返回 (pw_bytes|None, ok)。"""
    pw = simpledialog.askstring(
        "设置口令",
        "设置私钥保护口令\n（留空 = 不加密存储，不推荐）：",
        show='*', parent=root
    )
    if pw is None:
        return None, False
    if not pw:
        return None, True
    pw2 = simpledialog.askstring(
        "确认口令", "再次输入口令：", show='*', parent=root
    )
    if pw2 is None or pw != pw2:
        messagebox.showerror("错误", "两次输入的口令不一致。", parent=root)
        return None, False
    return pw.encode('utf-8'), True


def _ask_passphrase():
    """提示输入现有口令。返回 pw_bytes|None，或 False（取消）。"""
    pw = simpledialog.askstring(
        "输入口令", "私钥口令（留空表示无口令）：",
        show='*', parent=root
    )
    if pw is None:
        return False
    return pw.encode('utf-8') if pw else None


# -- 功能 -----------------------------------------------------------------------

def generate_key_pair():
    passphrase, ok = _ask_new_passphrase()
    if not ok:
        return
    enc = (serialization.BestAvailableEncryption(passphrase)
           if passphrase else serialization.NoEncryption())
    try:
        private_key = rsa.generate_private_key(
            public_exponent=65537, key_size=2048, backend=default_backend()
        )
        pri_pem = private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=enc
        )
        pub_pem = private_key.public_key().public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )
        with open(_MY_PRI, 'wb') as f:
            f.write(pri_pem)
        with open(_MY_PUB, 'wb') as f:
            f.write(pub_pem)
        note = " [口令保护]" if passphrase else " [警告：无口令保护]"
        _log(f"[OK] RSA-2048 密钥对已生成{note}。")
        _log(f"     私钥：{_MY_PRI}")
        _log(f"     公钥：{_MY_PUB}\n")
        _status("密钥对已生成。")
    except Exception as e:
        messagebox.showerror("错误", f"密钥生成失败：{e}", parent=root)


def import_public_key():
    path = filedialog.askopenfilename(
        title="选择对方公钥文件",
        filetypes=[("密钥文件", "*.pem *.txt"), ("所有文件", "*.*")],
        parent=root
    )
    if not path:
        return
    try:
        with open(path, 'rb') as f:
            data = f.read()
        serialization.load_pem_public_key(data, backend=default_backend())
        with open(_FR_PUB, 'wb') as f:
            f.write(data)
        _log("[OK] 对方公钥导入成功。\n")
        _status("对方公钥已导入。")
        import_btn.configure(text="[ 公钥已导入 ]")
    except Exception as e:
        messagebox.showerror("错误", f"导入失败：{e}", parent=root)


def encrypt_data():
    if os.path.isfile(_FR_PUB):
        key_path = _FR_PUB
    elif os.path.isfile(_MY_PUB):
        if not messagebox.askyesno(
            "未找到对方公钥",
            "未找到对方公钥。\n是否使用自己的公钥加密？",
            parent=root
        ):
            return
        key_path = _MY_PUB
    else:
        messagebox.showerror(
            "错误", "未找到公钥，请先生成密钥对。", parent=root
        )
        return
    try:
        with open(key_path, 'rb') as f:
            pub_key = serialization.load_pem_public_key(f.read(), backend=default_backend())
    except Exception as e:
        messagebox.showerror("错误", f"无法加载公钥：{e}", parent=root)
        return

    data = text_entry.get("1.0", tk.END).strip().encode('utf-8')
    if not data:
        _status("输入为空。")
        return
    try:
        # RSA-2048 + OAEP-SHA256：每块最大明文 190 字节
        MAX = 190
        encrypted = b''.join(
            pub_key.encrypt(
                data[i:i + MAX],
                padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                             algorithm=hashes.SHA256(), label=None)
            )
            for i in range(0, len(data), MAX)
        )
        hex_ct = encrypted.hex()
        pyperclip.copy(hex_ct)
        output_text.configure(state=tk.NORMAL)
        output_text.delete("1.0", tk.END)
        output_text.insert(tk.END, hex_ct)
        output_text.configure(state=tk.DISABLED)
        _status("加密完成 - 已自动复制到剪贴板。")
    except Exception as e:
        messagebox.showerror("错误", f"加密失败：{e}", parent=root)


def decrypt_data():
    if not os.path.isfile(_MY_PRI):
        messagebox.showerror(
            "错误", "未找到私钥，请先生成密钥对。", parent=root
        )
        return
    passphrase = _ask_passphrase()
    if passphrase is False:
        return
    try:
        with open(_MY_PRI, 'rb') as f:
            content = f.read()
        pri_key = serialization.load_pem_private_key(
            content, password=passphrase, backend=default_backend()
        )
    except Exception as e:
        messagebox.showerror("错误", f"无法加载私钥：{e}", parent=root)
        return

    hex_ct = text_entry.get("1.0", tk.END).strip()
    try:
        ciphertext = bytes.fromhex(hex_ct)
    except ValueError:
        messagebox.showerror("错误", "输入的十六进制密文无效。", parent=root)
        return
    try:
        # 每个 RSA-2048 密文块为 256 字节
        SEG = 256
        plaintext = b''.join(
            pri_key.decrypt(
                ciphertext[i:i + SEG],
                padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                             algorithm=hashes.SHA256(), label=None)
            )
            for i in range(0, len(ciphertext), SEG)
        )
        output_text.configure(state=tk.NORMAL)
        output_text.delete("1.0", tk.END)
        output_text.insert(tk.END, plaintext.decode('utf-8'))
        output_text.configure(state=tk.DISABLED)
        _status("解密完成。")
    except Exception as e:
        messagebox.showerror("错误", f"解密失败：{e}", parent=root)


def clear_all():
    text_entry.delete("1.0", tk.END)
    output_text.configure(state=tk.NORMAL)
    output_text.delete("1.0", tk.END)
    output_text.configure(state=tk.DISABLED)
    _status("已清空。")


def copy_output():
    content = output_text.get("1.0", tk.END).strip()
    if content:
        pyperclip.copy(content)
        _status("输出已复制到剪贴板。")


# -- 界面 -----------------------------------------------------------------------

root = tk.Tk()
root.title("CryphoCat")
root.geometry("600x420")
root.resizable(False, False)
root.configure(bg=_BG)

status_var = tk.StringVar(value="就绪。")


def _btn(parent, text, cmd, width=None):
    kw = dict(bg=_BFBG, fg=_FG, activebackground=_BFAC, activeforeground=_FG,
              font=_BFNT, relief=tk.RIDGE, bd=1, cursor='hand2')
    if width is not None:
        kw['width'] = width
    return tk.Button(parent, text=text, command=cmd, **kw)


def _lbl(parent, text):
    return tk.Label(parent, text=text, bg=_BG, fg=_FG, font=_SFNT)


# 顶部：密钥管理
top = tk.Frame(root, bg=_BG)
top.pack(fill=tk.X, padx=8, pady=(6, 2))
generate_btn = _btn(top, "[ 生成密钥对 ]", generate_key_pair, width=14)
generate_btn.pack(side=tk.LEFT, padx=(0, 4))
import_btn = _btn(top, "[ 导入对方公钥 ]", import_public_key, width=16)
import_btn.pack(side=tk.LEFT)
_lbl(top, "  RSA-2048 / OAEP-SHA256").pack(side=tk.RIGHT)

# 输入区域
in_f = tk.Frame(root, bg=_BG)
in_f.pack(fill=tk.BOTH, expand=True, padx=8, pady=2)
_lbl(in_f, "输入：").pack(anchor=tk.W)
text_entry = scrolledtext.ScrolledText(
    in_f, height=8, font=_FONT, bg='#111111', fg=_FG,
    insertbackground=_FG, relief=tk.FLAT, bd=1,
    selectbackground='#004400', selectforeground=_FG
)
text_entry.pack(fill=tk.BOTH, expand=True)

# 操作按钮行
act = tk.Frame(root, bg=_BG)
act.pack(fill=tk.X, padx=8, pady=4)
_btn(act, "[ 加密 ]", encrypt_data, width=10).pack(side=tk.LEFT, padx=(0, 4))
_btn(act, "[ 清空 ]", clear_all, width=8).pack(side=tk.LEFT, padx=(0, 4))
_btn(act, "[ 解密 ]", decrypt_data, width=10).pack(side=tk.LEFT, padx=(0, 4))
_btn(act, "[ 复制输出 ]", copy_output, width=12).pack(side=tk.RIGHT)

# 输出区域
out_f = tk.Frame(root, bg=_BG)
out_f.pack(fill=tk.BOTH, expand=True, padx=8, pady=(0, 2))
_lbl(out_f, "输出：").pack(anchor=tk.W)
output_text = scrolledtext.ScrolledText(
    out_f, height=6, font=_FONT, bg='#0a0a0a', fg=_OFGC,
    insertbackground=_FG, relief=tk.FLAT, bd=1,
    selectbackground='#004400', selectforeground=_FG,
    state=tk.DISABLED
)
output_text.pack(fill=tk.BOTH, expand=True)

# 状态栏
tk.Label(
    root, textvariable=status_var, bg='#0a0a0a', fg='#226622',
    font=_SFNT, anchor=tk.W
).pack(fill=tk.X, padx=8, pady=(0, 4))

root.mainloop()
