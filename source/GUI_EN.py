import os
import tkinter as tk
from tkinter import messagebox, scrolledtext, simpledialog
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.backends import default_backend
import tkinter.filedialog as filedialog
import pyperclip

# -- Key storage ---------------------------------------------------------------
_BASE   = 'RSAkeys'
_MY_DIR = os.path.join(_BASE, 'my')
_FR_DIR = os.path.join(_BASE, 'friend')
_MY_PRI = os.path.join(_MY_DIR, 'private.pem')
_MY_PUB = os.path.join(_MY_DIR, 'public.pem')
_FR_PUB = os.path.join(_FR_DIR, 'public.pem')

for _d in (_MY_DIR, _FR_DIR):
    os.makedirs(_d, exist_ok=True)

# -- Retro terminal theme ------------------------------------------------------
_BG   = '#0d0d0d'
_FG   = '#33cc33'
_OFGC = '#22aa22'
_BFBG = '#1c1c1c'
_BFAC = '#003300'
_FONT = ('Courier', 9)
_BFNT = ('Courier', 9, 'bold')
_SFNT = ('Courier', 8)


# -- Helpers -------------------------------------------------------------------

def _log(msg):
    output_text.configure(state=tk.NORMAL)
    output_text.insert(tk.END, msg + '\n')
    output_text.see(tk.END)
    output_text.configure(state=tk.DISABLED)


def _status(msg):
    status_var.set(msg)


def _ask_new_passphrase():
    """Prompt for new passphrase with confirmation. Returns (pw_bytes|None, ok)."""
    pw = simpledialog.askstring(
        "Set Passphrase",
        "Passphrase for private key\n(blank = no encryption, NOT recommended):",
        show='*', parent=root
    )
    if pw is None:
        return None, False
    if not pw:
        return None, True
    pw2 = simpledialog.askstring(
        "Confirm Passphrase", "Confirm passphrase:", show='*', parent=root
    )
    if pw2 is None or pw != pw2:
        messagebox.showerror("Error", "Passphrases do not match.", parent=root)
        return None, False
    return pw.encode('utf-8'), True


def _ask_passphrase():
    """Prompt for existing passphrase. Returns pw_bytes|None, or False to abort."""
    pw = simpledialog.askstring(
        "Passphrase", "Private key passphrase (blank if none):",
        show='*', parent=root
    )
    if pw is None:
        return False
    return pw.encode('utf-8') if pw else None


# -- Actions -------------------------------------------------------------------

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
        note = " [passphrase protected]" if passphrase else " [WARNING: no passphrase]"
        _log(f"[OK] RSA-2048 key pair generated{note}.")
        _log(f"     Private : {_MY_PRI}")
        _log(f"     Public  : {_MY_PUB}\n")
        _status("Keys generated.")
    except Exception as e:
        messagebox.showerror("Error", f"Key generation failed: {e}", parent=root)


def import_public_key():
    path = filedialog.askopenfilename(
        title="Select friend's public key",
        filetypes=[("Key files", "*.pem *.txt"), ("All files", "*.*")],
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
        _log("[OK] Friend's public key imported.\n")
        _status("Friend's key imported.")
        import_btn.configure(text="[ Key Imported ]")
    except Exception as e:
        messagebox.showerror("Error", f"Import failed: {e}", parent=root)


def encrypt_data():
    if os.path.isfile(_FR_PUB):
        key_path = _FR_PUB
    elif os.path.isfile(_MY_PUB):
        if not messagebox.askyesno(
            "No friend key",
            "Friend's public key not found.\nUse your own public key?",
            parent=root
        ):
            return
        key_path = _MY_PUB
    else:
        messagebox.showerror(
            "Error", "No public key found. Generate a key pair first.", parent=root
        )
        return
    try:
        with open(key_path, 'rb') as f:
            pub_key = serialization.load_pem_public_key(f.read(), backend=default_backend())
    except Exception as e:
        messagebox.showerror("Error", f"Cannot load public key: {e}", parent=root)
        return

    data = text_entry.get("1.0", tk.END).strip().encode('utf-8')
    if not data:
        _status("Input is empty.")
        return
    try:
        # RSA-2048 + OAEP-SHA256: max plaintext per block = 190 bytes
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
        _status("Encrypted - copied to clipboard.")
    except Exception as e:
        messagebox.showerror("Error", f"Encryption failed: {e}", parent=root)


def decrypt_data():
    if not os.path.isfile(_MY_PRI):
        messagebox.showerror(
            "Error", "Private key not found. Generate a key pair first.", parent=root
        )
        return
    passphrase = _ask_passphrase()
    if passphrase is False:
        return
    try:
        with open(_MY_PRI, 'rb') as f:
            pri_key = serialization.load_pem_private_key(
                f.read(), password=passphrase, backend=default_backend()
            )
    except Exception as e:
        messagebox.showerror("Error", f"Cannot load private key: {e}", parent=root)
        return

    hex_ct = text_entry.get("1.0", tk.END).strip()
    try:
        ciphertext = bytes.fromhex(hex_ct)
    except ValueError:
        messagebox.showerror("Error", "Input is not valid hex ciphertext.", parent=root)
        return
    try:
        # Each RSA-2048 ciphertext block is 256 bytes
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
        _status("Decrypted.")
    except Exception as e:
        messagebox.showerror("Error", f"Decryption failed: {e}", parent=root)


def clear_all():
    text_entry.delete("1.0", tk.END)
    output_text.configure(state=tk.NORMAL)
    output_text.delete("1.0", tk.END)
    output_text.configure(state=tk.DISABLED)
    _status("Cleared.")


def copy_output():
    content = output_text.get("1.0", tk.END).strip()
    if content:
        pyperclip.copy(content)
        _status("Output copied to clipboard.")


# -- GUI -----------------------------------------------------------------------

root = tk.Tk()
root.title("CryphoCat")
root.geometry("600x420")
root.resizable(False, False)
root.configure(bg=_BG)

status_var = tk.StringVar(value="Ready.")


def _btn(parent, text, cmd, width=None):
    kw = dict(bg=_BFBG, fg=_FG, activebackground=_BFAC, activeforeground=_FG,
              font=_BFNT, relief=tk.RIDGE, bd=1, cursor='hand2')
    if width is not None:
        kw['width'] = width
    return tk.Button(parent, text=text, command=cmd, **kw)


def _lbl(parent, text):
    return tk.Label(parent, text=text, bg=_BG, fg=_FG, font=_SFNT)


# Top row: key management
top = tk.Frame(root, bg=_BG)
top.pack(fill=tk.X, padx=8, pady=(6, 2))
generate_btn = _btn(top, "[ Gen Keys ]", generate_key_pair, width=14)
generate_btn.pack(side=tk.LEFT, padx=(0, 4))
import_btn = _btn(top, "[ Import Friend Key ]", import_public_key, width=20)
import_btn.pack(side=tk.LEFT)
_lbl(top, "  RSA-2048 / OAEP-SHA256").pack(side=tk.RIGHT)

# Input area
in_f = tk.Frame(root, bg=_BG)
in_f.pack(fill=tk.BOTH, expand=True, padx=8, pady=2)
_lbl(in_f, "INPUT:").pack(anchor=tk.W)
text_entry = scrolledtext.ScrolledText(
    in_f, height=8, font=_FONT, bg='#111111', fg=_FG,
    insertbackground=_FG, relief=tk.FLAT, bd=1,
    selectbackground='#004400', selectforeground=_FG
)
text_entry.pack(fill=tk.BOTH, expand=True)

# Action row
act = tk.Frame(root, bg=_BG)
act.pack(fill=tk.X, padx=8, pady=4)
_btn(act, "[ Encrypt ]", encrypt_data, width=12).pack(side=tk.LEFT, padx=(0, 4))
_btn(act, "[ Clear ]", clear_all, width=8).pack(side=tk.LEFT, padx=(0, 4))
_btn(act, "[ Decrypt ]", decrypt_data, width=12).pack(side=tk.LEFT, padx=(0, 4))
_btn(act, "[ Copy Output ]", copy_output, width=14).pack(side=tk.RIGHT)

# Output area
out_f = tk.Frame(root, bg=_BG)
out_f.pack(fill=tk.BOTH, expand=True, padx=8, pady=(0, 2))
_lbl(out_f, "OUTPUT:").pack(anchor=tk.W)
output_text = scrolledtext.ScrolledText(
    out_f, height=6, font=_FONT, bg='#0a0a0a', fg=_OFGC,
    insertbackground=_FG, relief=tk.FLAT, bd=1,
    selectbackground='#004400', selectforeground=_FG,
    state=tk.DISABLED
)
output_text.pack(fill=tk.BOTH, expand=True)

# Status bar
tk.Label(
    root, textvariable=status_var, bg='#0a0a0a', fg='#226622',
    font=_SFNT, anchor=tk.W
).pack(fill=tk.X, padx=8, pady=(0, 4))

root.mainloop()
