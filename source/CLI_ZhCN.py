import os
import time
import hashlib
import getpass
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.backends import default_backend
import pyperclip
import colorama
from colorama import Fore, Style

colorama.init()

_BASE   = 'RSAkeys'
_MY_DIR = os.path.join(_BASE, 'my')
_FR_DIR = os.path.join(_BASE, 'friend')
_MY_PRI = os.path.join(_MY_DIR, 'private.pem')
_MY_PUB = os.path.join(_MY_DIR, 'public.pem')
_FR_PUB = os.path.join(_FR_DIR, 'public.pem')

for _d in (_MY_DIR, _FR_DIR):
    os.makedirs(_d, exist_ok=True)


def generate_key_pair():
    raw = input('密钥长度（1024/2048/4096，默认 4096）：').strip()
    if not raw:
        key_size = 4096
    else:
        try:
            key_size = int(raw)
            if key_size not in (1024, 2048, 4096):
                print('请从 1024、2048、4096 中选择一个。')
                return
        except ValueError:
            print('请输入数字。')
            return

    pw = getpass.getpass('设置私钥口令（留空 = 不加密存储，不推荐）：')
    if pw:
        pw2 = getpass.getpass('确认口令：')
        if pw != pw2:
            print('两次输入的口令不一致。')
            return
        enc = serialization.BestAvailableEncryption(pw.encode('utf-8'))
    else:
        enc = serialization.NoEncryption()
        print('[警告] 私钥将以明文形式存储，安全性较低。')

    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=key_size,
        backend=default_backend()
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
    print(f'RSA-{key_size} 密钥对已生成。')
    print(f'  私钥：{_MY_PRI}')
    print(f'  公钥：{_MY_PUB}')


def import_public_key():
    path = input('请输入对方公钥文件路径：').strip()
    if os.path.isfile(path):
        try:
            with open(path, 'rb') as f:
                data = f.read()
            serialization.load_pem_public_key(data, backend=default_backend())
            with open(_FR_PUB, 'wb') as f:
                f.write(data)
            print('公钥导入成功。')
        except Exception as e:
            print(f'导入失败：{e}')
    else:
        print('文件不存在，请检查路径。')


def encrypt_data():
    if os.path.isfile(_FR_PUB):
        key_path = _FR_PUB
    elif os.path.isfile(_MY_PUB):
        ans = input('未找到对方公钥，是否使用自己的公钥加密？(y/n)：').strip().lower()
        if ans != 'y':
            print('已取消。请先导入对方公钥或生成密钥对。')
            return
        key_path = _MY_PUB
    else:
        print('未找到公钥，请先生成密钥对。')
        return

    try:
        with open(key_path, 'rb') as f:
            pub_key = serialization.load_pem_public_key(f.read(), backend=default_backend())
    except Exception as e:
        print(f'无法加载公钥：{e}')
        return

    data = input('请输入要加密的文本：').encode('utf-8')
    if not data:
        print('输入为空。')
        return
    try:
        MAX = 190  # RSA-2048 + OAEP-SHA256
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
        print('加密结果（已复制到剪贴板）：', hex_ct)
    except Exception as e:
        print(f'加密失败：{e}')


def decrypt_data():
    if not os.path.isfile(_MY_PRI):
        print('未找到私钥，请先生成密钥对。')
        return
    try:
        raw_pw = getpass.getpass('私钥口令（无口令则留空）：')
        pw_bytes = raw_pw.encode('utf-8') if raw_pw else None
        with open(_MY_PRI, 'rb') as f:
            pri_key = serialization.load_pem_private_key(
                f.read(), password=pw_bytes, backend=default_backend()
            )
    except Exception as e:
        print(f'无法加载私钥：{e}')
        return

    hex_ct = input('请输入十六进制密文：').strip()
    try:
        ciphertext = bytes.fromhex(hex_ct)
    except ValueError:
        print('十六进制密文无效。')
        return
    try:
        SEG = 256
        plaintext = b''.join(
            pri_key.decrypt(
                ciphertext[i:i + SEG],
                padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                             algorithm=hashes.SHA256(), label=None)
            )
            for i in range(0, len(ciphertext), SEG)
        )
        print('解密结果：', plaintext.decode('utf-8'))
    except Exception as e:
        print(f'解密失败：{e}')


def calculate_hashes():
    path = input('请输入文件路径：').strip()
    if os.path.isfile(path):
        with open(path, 'rb') as f:
            data = f.read()
        print('MD5   :', hashlib.md5(data).hexdigest())
        print('SHA1  :', hashlib.sha1(data).hexdigest())
        print('SHA256:', hashlib.sha256(data).hexdigest())
    else:
        print('文件不存在，请检查路径。')


def compare_hashes():
    path = input('请输入文件路径：').strip()
    file_hash = input('请输入哈希值（MD5/SHA1/SHA256）：').strip()
    if not os.path.isfile(path):
        print('文件不存在，请检查路径。')
        return
    if not file_hash:
        print('未输入哈希值。')
        return
    with open(path, 'rb') as f:
        data = f.read()
    checks = {
        'MD5':    hashlib.md5(data).hexdigest(),
        'SHA1':   hashlib.sha1(data).hexdigest(),
        'SHA256': hashlib.sha256(data).hexdigest(),
    }
    matched = False
    for algo, digest in checks.items():
        if file_hash.lower() == digest:
            print(Fore.GREEN + f'{algo} 校验通过：{path}' + Style.RESET_ALL)
            print(Fore.GREEN + f'哈希值：{digest}' + Style.RESET_ALL)
            matched = True
    if not matched:
        print(Fore.RED + '哈希校验失败 —— 文件可能已被篡改。' + Style.RESET_ALL)


def main_menu():
    print("----------------")
    print("1. 生成 RSA 密钥对")
    print("2. 导入对方公钥")
    print("3. 加密文本")
    print("4. 解密文本")
    print("5. 计算文件哈希")
    print("6. 校验文件哈希")
    print("7. 退出")
    print("----------------")


while True:
    main_menu()
    choice = input('请输入选项：').strip()
    if choice == '1':
        generate_key_pair()
        time.sleep(1)
    elif choice == '2':
        import_public_key()
    elif choice == '3':
        encrypt_data()
        input('按下回车继续...')
    elif choice == '4':
        decrypt_data()
        input('按下回车继续...')
    elif choice == '5':
        calculate_hashes()
        input('按下回车继续...')
    elif choice == '6':
        compare_hashes()
        input('按下回车继续...')
    elif choice == '7':
        print('再见。')
        break
    else:
        print('选项无效，请重新输入。')
        input('按下回车继续...')
