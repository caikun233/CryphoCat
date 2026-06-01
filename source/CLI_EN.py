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
    raw = input('Key length (1024/2048/4096, default 4096): ').strip()
    if not raw:
        key_size = 4096
    else:
        try:
            key_size = int(raw)
            if key_size not in (1024, 2048, 4096):
                print('Choose one of: 1024, 2048, 4096')
                return
        except ValueError:
            print('Please enter a number.')
            return

    pw = getpass.getpass('Set passphrase (blank = no encryption, NOT recommended): ')
    if pw:
        pw2 = getpass.getpass('Confirm passphrase: ')
        if pw != pw2:
            print('Passphrases do not match.')
            return
        enc = serialization.BestAvailableEncryption(pw.encode('utf-8'))
    else:
        enc = serialization.NoEncryption()
        print('[WARNING] Private key will be stored without a passphrase.')

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
    print(f'RSA-{key_size} key pair generated.')
    print(f'  Private: {_MY_PRI}')
    print(f'  Public : {_MY_PUB}')


def import_public_key():
    path = input('Path to friend\'s public key: ').strip()
    if os.path.isfile(path):
        try:
            with open(path, 'rb') as f:
                data = f.read()
            serialization.load_pem_public_key(data, backend=default_backend())
            with open(_FR_PUB, 'wb') as f:
                f.write(data)
            print('Public key imported successfully.')
        except Exception as e:
            print(f'Import failed: {e}')
    else:
        print('File not found. Check the path.')


def encrypt_data():
    if os.path.isfile(_FR_PUB):
        key_path = _FR_PUB
    elif os.path.isfile(_MY_PUB):
        ans = input('Friend\'s key not found. Use your own public key? (y/n): ').strip().lower()
        if ans != 'y':
            print('Cancelled. Import a friend\'s key or generate your key pair first.')
            return
        key_path = _MY_PUB
    else:
        print('No public key available. Generate a key pair first.')
        return

    try:
        with open(key_path, 'rb') as f:
            pub_key = serialization.load_pem_public_key(f.read(), backend=default_backend())
    except Exception as e:
        print(f'Cannot load public key: {e}')
        return

    data = input('Text to encrypt: ').encode('utf-8')
    if not data:
        print('Nothing to encrypt.')
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
        print('Encrypted (copied to clipboard):', hex_ct)
    except Exception as e:
        print(f'Encryption failed: {e}')


def decrypt_data():
    if not os.path.isfile(_MY_PRI):
        print('Private key not found. Generate a key pair first.')
        return
    try:
        raw_pw = getpass.getpass('Private key passphrase (blank if none): ')
        pw_bytes = raw_pw.encode('utf-8') if raw_pw else None
        with open(_MY_PRI, 'rb') as f:
            pri_key = serialization.load_pem_private_key(
                f.read(), password=pw_bytes, backend=default_backend()
            )
    except Exception as e:
        print(f'Cannot load private key: {e}')
        return

    hex_ct = input('Ciphertext (hex): ').strip()
    try:
        ciphertext = bytes.fromhex(hex_ct)
    except ValueError:
        print('Invalid hex ciphertext.')
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
        print('Decrypted:', plaintext.decode('utf-8'))
    except Exception as e:
        print(f'Decryption failed: {e}')


def calculate_hashes():
    path = input('File path: ').strip()
    if os.path.isfile(path):
        with open(path, 'rb') as f:
            data = f.read()
        print('MD5   :', hashlib.md5(data).hexdigest())
        print('SHA1  :', hashlib.sha1(data).hexdigest())
        print('SHA256:', hashlib.sha256(data).hexdigest())
    else:
        print('File not found. Check the path.')


def compare_hashes():
    path = input('File path: ').strip()
    file_hash = input('Expected hash (MD5/SHA1/SHA256): ').strip()
    if not os.path.isfile(path):
        print('File not found. Check the path.')
        return
    if not file_hash:
        print('No hash provided.')
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
            print(Fore.GREEN + f'{algo} MATCH: {path}' + Style.RESET_ALL)
            print(Fore.GREEN + f'Value : {digest}' + Style.RESET_ALL)
            matched = True
    if not matched:
        print(Fore.RED + 'HASH CHECK FAILED – file may have been tampered with.' + Style.RESET_ALL)


def main_menu():
    print("----------------")
    print("1. Generate RSA key pair")
    print("2. Import friend\'s public key")
    print("3. Encrypt text")
    print("4. Decrypt text")
    print("5. Calculate file hashes")
    print("6. Verify file hash")
    print("7. Exit")
    print("----------------")


while True:
    main_menu()
    choice = input('Choice: ').strip()
    if choice == '1':
        generate_key_pair()
        time.sleep(1)
    elif choice == '2':
        import_public_key()
    elif choice == '3':
        encrypt_data()
        input('Press Enter to continue...')
    elif choice == '4':
        decrypt_data()
        input('Press Enter to continue...')
    elif choice == '5':
        calculate_hashes()
        input('Press Enter to continue...')
    elif choice == '6':
        compare_hashes()
        input('Press Enter to continue...')
    elif choice == '7':
        print('Goodbye.')
        break
    else:
        print('Invalid choice.')
        input('Press Enter to continue...')
