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
    print('no path found like RSAKeys, creating\n')
    os.makedirs('RSAkeys')

if not os.path.exists('RSAkeys/my'):
    os.makedirs('RSAkeys/my')

if not os.path.exists('RSAkeys/friend'):
    os.makedirs('RSAkeys/friend')

def generate_key_pair():
    input_key_size = input('input key length(1024/2048/4096 default):')
    if not input_key_size:
        input_key_size = 4096
    else:
        try:
            input_key_size = int(input_key_size)
            if input_key_size not in [1024, 2048, 4096]:
                print('you can only choose one in 1024/2048/4096')
                return
        except ValueError:
            print('are you inputing numbers?')
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

    print('RSA key pair generated, saved at RSAkeys/my/public.txt and RSAkeys/my/private.txt')

def import_public_key():
    path = input('input path of public_key you want import：')
    if os.path.isfile(path):
        with open(path, 'rb') as f:
            public_pem = f.read()
            with open('RSAkeys/friend/public.txt', 'wb') as fw:
                fw.write(public_pem)
        print('public key import success!')
    else:
        print('not found, check the path')

def encrypt_data():
    if os.path.exists('RSAkeys/friend/public.txt'):
        with open('RSAkeys/friend/public.txt', 'rb') as f:
            public_pem = f.read()
            public_key = serialization.load_pem_public_key(public_pem, backend=default_backend())

            data = input('input the text you want encrypt: ').encode('utf-8')
            ciphertext = public_key.encrypt(data, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                               algorithm=hashes.SHA256(), label=None))
            print('hex text after encrypting(auto copied): ', ciphertext.hex())
    else:
        use_own_public_key = input('no imported public_key，Use the key generated? (T=True，F=False): ')
        if use_own_public_key == 'T':
            with open('RSAkeys/my/public.txt', 'rb') as f:
                public_pem = f.read()
                public_key = serialization.load_pem_public_key(public_pem, backend=default_backend())

                data = input('input the text you want encrypt: ').encode('utf-8')
                ciphertext = public_key.encrypt(data, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                                   algorithm=hashes.SHA256(), label=None))
                print('hex text after encrypting(auto copied): ', ciphertext.hex())
        else:
            print('plz import others public_key or generate your key pair')
    pyperclip.copy(ciphertext.hex())

def decrypt_data():
    with open('RSAkeys/my/private.txt', 'rb') as f:
        private_pem = f.read()
        private_key = serialization.load_pem_private_key(private_pem, password=None, backend=default_backend())

        ciphertext = bytes.fromhex(input('input the text you want decrypt: '))

        data = private_key.decrypt(ciphertext, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                            algorithm=hashes.SHA256(), label=None))
        print('text after decrypting: ', data.decode('utf-8'))

def calculate_hashes():
    path = input('input file path: ')
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
        print('not found, check path')

def compare_hashes():
    path = input('input file path:')
    file_hash = input('input HEX hash(MD5/SHA1/SHA256): ')
    if os.path.isfile(path):
        if file_hash:
            with open(path, 'rb') as f:
                data = f.read()
                
                md5_hash = hashlib.md5(data).hexdigest()
                sha1_hash = hashlib.sha1(data).hexdigest()
                sha256_hash = hashlib.sha256(data).hexdigest()
                
                match_found = False
                if file_hash == md5_hash:
                    print(Fore.GREEN + 'MD5 MATCH FOUND:' + Style.RESET_ALL, path)
                    print(Fore.GREEN + 'with value:' + Style.RESET_ALL, md5_hash)
                    match_found = True
                if file_hash == sha1_hash:
                    print(Fore.GREEN + 'SHA1 MATCH FOUND:' + Style.RESET_ALL, path)
                    print(Fore.GREEN + 'with value:' + Style.RESET_ALL, sha1_hash)
                    match_found = True
                if file_hash == sha256_hash:
                    print(Fore.GREEN + 'SHA256 MATCH FOUND:' + Style.RESET_ALL, path)
                    print(Fore.GREEN + 'with value:' + Style.RESET_ALL, sha256_hash)
                    match_found = True
                if not match_found:
                    print(Fore.RED + 'HASH CHECK FAILED, FILE MAY BE TAMPERED' + Style.RESET_ALL)
        else:
            print('you really inputted hash?')
    else:
        print('file not found, check path')

def main_menu():
    print("----------------")
    print("1. Generate RSA key pair.")
    print("2. Import public key from others.")
    print("3. Encrypt data.")
    print("4. Decrypt data.")
    print("5. Calculate file hash")
    print("6. Check file hash")
    print("7. Say Goodbye")
    print("----------------")

while True:
    main_menu()
    choice = input('what is your choice：')

    if choice == '1':
        generate_key_pair()
        time.sleep(1)
    elif choice == '2':
        import_public_key()
    elif choice == '3':
        encrypt_data()
        input("press enter to continue...")
    elif choice == '4':
        decrypt_data()
        input("press enter to continue...")
    elif choice == '5':
        calculate_hashes()
        input("press enter to continue...")
    elif choice == '6':
        compare_hashes()
        input("press enter to continue...")
    elif choice == '7':
        print('Never gonna tell↓ a↑ liiiie↓ and hurt you.')
        break
    else:
        print('? what did you just press')
        input("press enter to continue...")
