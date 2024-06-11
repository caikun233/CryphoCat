import os
import time
import hashlib
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.backends import default_backend
import pyperclip
import colorama
from colorama import Fore, Style

# 初始化 colorama
colorama.init()

current_dir = os.path.dirname(os.path.abspath(__file__))

if not os.path.exists('RSAkeys'):
    print('检测到没有RSAkeys文件夹，自动创建。\n')
    os.makedirs('RSAkeys')

if not os.path.exists('RSAkeys/my'):
    os.makedirs('RSAkeys/my')

if not os.path.exists('RSAkeys/friend'):
    os.makedirs('RSAkeys/friend')

def generate_key_pair():
    input_key_size = input('输入密钥长度(1024/2048/默认4096):')
    if not input_key_size:
        input_key_size = 4096
    else:
        try:
            input_key_size = int(input_key_size)
            if input_key_size not in [1024, 2048, 4096]:
                print('三个里面选一个，别的不认')
                return
        except ValueError:
            print('？你输的是数字？')
            return
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=input_key_size,
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

    print('RSA密钥对生成成功，并已保存到RSAkeys/my/public.txt和RSAkeys/my/private.txt')

def import_public_key():
    path = input('请输入要导入的公钥文件的路径：')
    if os.path.isfile(path):
        with open(path, 'rb') as f:
            public_pem = f.read()
            with open('RSAkeys/friend/public.txt', 'wb') as fw:
                fw.write(public_pem)
        print('公钥导入成功！')
    else:
        print('文件不存在，请检查路径是否正确。')

def encrypt_data():
    if os.path.exists('RSAkeys/friend/public.txt'):
        with open('RSAkeys/friend/public.txt', 'rb') as f:
            public_pem = f.read()
            public_key = serialization.load_pem_public_key(public_pem, backend=default_backend())

            data = input('请输入要加密的数据：').encode('utf-8')
            ciphertext = public_key.encrypt(data, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                               algorithm=hashes.SHA256(), label=None))
            print('加密后的数据为：', ciphertext.hex())
    else:
        use_own_public_key = input('未导入他人公钥，是否使用自己的公钥进行加密？ (1=是，2=否): ')
        if use_own_public_key == '1':
            with open('RSAkeys/my/public.txt', 'rb') as f:
                public_pem = f.read()
                public_key = serialization.load_pem_public_key(public_pem, backend=default_backend())

                data = input('请输入要加密的数据：').encode('utf-8')
                ciphertext = public_key.encrypt(data, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                                   algorithm=hashes.SHA256(), label=None))
                print('加密后的数据为（已经自动复制）：', ciphertext.hex())
        else:
            print('请先导入他人公钥或生成自己的密钥对。')
    pyperclip.copy(ciphertext.hex())

def decrypt_data():
    with open('RSAkeys/my/private.txt', 'rb') as f:
        private_pem = f.read()
        private_key = serialization.load_pem_private_key(private_pem, password=None, backend=default_backend())

        ciphertext = bytes.fromhex(input('请输入要解密的数据：'))

        data = private_key.decrypt(ciphertext, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                            algorithm=hashes.SHA256(), label=None))
        print('解密后的数据为：', data.decode('utf-8'))

def calculate_hashes():
    path = input('请输入文件路径：')
    if os.path.isfile(path):
        with open(path, 'rb') as f:
            data = f.read()
            
            md5_hash = hashlib.md5(data).hexdigest()
            sha1_hash = hashlib.sha1(data).hexdigest()
            sha256_hash = hashlib.sha256(data).hexdigest()

            print('MD5:', md5_hash)
            print('SHA1:', sha1_hash)
            print('SHA256:', sha256_hash)
    else:
        print('文件不存在，请检查路径是否正确。')

def compare_hashes():
    path = input('请输入文件路径：')
    file_hash = input('请输入哈希（MD5/SHA1/SHA256）：')
    if os.path.isfile(path):
        if file_hash:
            with open(path, 'rb') as f:
                data = f.read()
                
                md5_hash = hashlib.md5(data).hexdigest()
                sha1_hash = hashlib.sha1(data).hexdigest()
                sha256_hash = hashlib.sha256(data).hexdigest()
                
                match_found = False
                if file_hash == md5_hash:
                    print(Fore.GREEN + 'MD5校验通过:' + Style.RESET_ALL, path)
                    print(Fore.GREEN + '哈希值:' + Style.RESET_ALL, md5_hash)
                    match_found = True
                if file_hash == sha1_hash:
                    print(Fore.GREEN + 'SHA1校验通过:' + Style.RESET_ALL, path)
                    print(Fore.GREEN + '哈希值:' + Style.RESET_ALL, sha1_hash)
                    match_found = True
                if file_hash == sha256_hash:
                    print(Fore.GREEN + 'SHA256校验通过:' + Style.RESET_ALL, path)
                    print(Fore.GREEN + '哈希值:' + Style.RESET_ALL, sha256_hash)
                    match_found = True
                if not match_found:
                    print(Fore.RED + '哈希校验失败，文件可能已被修改' + Style.RESET_ALL)
        else:
            print('输哈希了吗你就点啊？')
    else:
        print('文件不存在，请检查路径是否正确。')

def main_menu():
    print("----------------")
    print("1. 生成RSA密钥对")
    print("2. 导入他人公钥")
    print("3. 加密数据")
    print("4. 解密数据")
    print("5. 计算文件哈希值")
    print("6. 校验哈希")
    print("7. 退出")
    print("----------------")

while True:
    main_menu()
    choice = input('请输入选项：')

    if choice == '1':
        generate_key_pair()
        time.sleep(1)
    elif choice == '2':
        import_public_key()
    elif choice == '3':
        encrypt_data()
        input("按下回车继续...")
    elif choice == '4':
        decrypt_data()
        input("按下回车继续...")
    elif choice == '5':
        calculate_hashes()
        input("按下回车继续...")
    elif choice == '6':
        compare_hashes()
        input("按下回车继续...")
    elif choice == '7':
        print('程序已退出。')
        break
    else:
        print('选项不正确，请重新选择。')
        input("按下回车继续...")
