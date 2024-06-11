import os
import time
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.backends import default_backend
import pyperclip

current_dir = os.path.dirname(os.path.abspath(__file__))

if not os.path.exists('RSAkeys'):
    print('Detected no `RSAkeys` folder, creating it automatically.\n')
    os.makedirs('RSAkeys')

if not os.path.exists('RSAkeys/my'):
    os.makedirs('RSAkeys/my')

if not os.path.exists('RSAkeys/friend'):
    os.makedirs('RSAkeys/friend')

def generate_key_pair():
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=4096,
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

    print('RSA key pair generated successfully! Saved to RSAkeys/my/public.txt and RSAkeys/my/private.txt')

def import_public_key():
    path = input('Input the path to pubkey file you want to import: ')
    if os.path.isfile(path):
        with open(path, 'rb') as f:
            public_pem = f.read()
            with open('RSAkeys/friend/public.txt', 'wb') as fw:
                fw.write(public_pem)
        print('SUCCESS!')
    else:
        print('NO FILE DETECTED')

def encrypt_data():
    if os.path.exists('RSAkeys/friend/public.txt'):
        with open('RSAkeys/friend/public.txt', 'rb') as f:
            public_pem = f.read()
            public_key = serialization.load_pem_public_key(public_pem, backend=default_backend())

            data = input('Input data you want ro encrypt:').encode('utf-8')
            ciphertext = public_key.encrypt(data, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                               algorithm=hashes.SHA256(), label=None))
            print("Here's your encrypted data: (auto copied)", ciphertext.hex())
    else:
        use_own_public_key = input("Your friend's public key was not found. Do you want to use your own public key for encryption? (1=Yes，2=No): ")
        if use_own_public_key == '1':
            with open('RSAkeys/my/public.txt', 'rb') as f:
                public_pem = f.read()
                public_key = serialization.load_pem_public_key(public_pem, backend=default_backend())

                data = input('Input data you want ro encrypt:').encode('utf-8')
                ciphertext = public_key.encrypt(data, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                                   algorithm=hashes.SHA256(), label=None))
                print("Here's your encrypted data: (auto copied)：", ciphertext.hex())
        else:
            print("You didn't import or generate any publickey.")
    pyperclip.copy(ciphertext.hex())

def decrypt_data():
    with open('RSAkeys/my/private.txt', 'rb') as f:
        private_pem = f.read()
        private_key = serialization.load_pem_private_key(private_pem, password=None, backend=default_backend())

        ciphertext = bytes.fromhex(input('Input data you want DECRYPT:'))

        data = private_key.decrypt(ciphertext, padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()),
                                                            algorithm=hashes.SHA256(), label=None))
        print('Here is your data decrypted: ', data.decode('utf-8'))

def main_menu():
    print("----------------")
    print("1. Generate RSA key pair.")
    print("2. Import public key from others.")
    print("3. Encrypt data.")
    print("4. Decrypt data.")
    print("5. Say Goodbye.")
    print("----------------")

while True:
    main_menu()
    choice = input('So what is your choice?')

    if choice == '1':
        generate_key_pair()
        time.sleep(3)
    elif choice == '2':
        import_public_key()
    elif choice == '3':
        encrypt_data()
        input("Press ENTER continue.")
    elif choice == '4':
        decrypt_data()
        input("Press ENTER continue.")
    elif choice == '5':
        print('Never gonna tell↓ a↑ liiiie↓ and hurt you.')
        break
    else:
        print('? what did you just press')
        input("Press ENTER continue.")
