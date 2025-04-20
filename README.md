# CryphoCat
~~***THE WHOLE PROJECT IS UNDER ALPHA DEVELOPING***~~

~~***THE WHOLE PROJECT IS UNDER ALPHA DEVELOPING***~~

~~***THE WHOLE PROJECT IS UNDER ALPHA DEVELOPING***~~

***THIS PROJECT WILL REMAKE WITH OTHER GUI FRAMEWORK LIKE QT SOON***

Encrypt your chatting with me, meow~

Currently only Simplified Chinese and English are supported.

![](https://img.shields.io/badge/python-v3.10-blue)
![GitHub](https://img.shields.io/github/license/caikun233/CryphoCat)

**中文README请参照**[README.zhcn.md](https://github.com/caikun233/CryphoCat/blob/main/README.zhcn.md)

* Support asymmetric encryption of arbitrary length text, currently using RSA2048 (I want to make it optional, but the development energy is not here).
* The key directory generation rules in the "two-person chat" scenario have been built in. You can Pull Request to add more.
* All offline processing, and open source, do not upload any data, decoupled from the chat software.
* As long as your chat software can guarantee that the information sent and received is the same, then the content of your communication will never be revealed.
* **Not Support Audio/Video At All, Unless there is a way to convert audio/video to text**.
* You can base64 encode the picture and send it out, and I may add the base64 encoding function to the software later. ~~Only idoits send images like this, just telling others what you want to say and don't use the xxxxing images.~~

## INSTALLATION & USAGE

### Linux / Windows CLI🔨

 1. I coded over Python3.10.9 & Windows 10 x64. You'd better install Python3.10 First.

 2. Just download CLI.py and 

    ```
    python CLI.py
    ```

 3. The program will automatically create a folder named "RSAkeys" in the same level directory, and create a folder of "my" and "friend" in it to distinguish the key between you and the chat partner.

 4. ~~The program only support zh-CN till y2023/m07/d01. I will upload EN version soon, so no more explanation here.~~ I uploaded EN version.

###  Windows GUI 🔨

1. **The GUI beta version completed.** 
2. No more EXE file in releases! Download from Actions! Thanks GitHub!**No x86 support, x64 only**.
3. Double click exe file, it looks ugly, right? I am not good at any art, but I will try my best to make it more beautiful.
4. The first text entry box is your friend's public key's path you want to input. The 2nd text entry box is where you input words to encrypt or decrypt.
5. ~~Also, The program only support zh-CN till y2023/m07/d01.~~ I uploaded EN version.

## Development Plan

- [x] Add English support for releases and source code.
- [x] Finish GUI version development.
- [x] Make RSA key length optional.(only CLI completed)
- [ ] Make GUI more beautiful.
- [ ] Try to add images base64 encode.
- [x] Let the encrypted text be copied automatically.
- [x] According to a Twitterer, we can add a function to calculate some hash of images to avoid the platform adding tracking watermarks to the image, also we can check them.(Now supports any file)
