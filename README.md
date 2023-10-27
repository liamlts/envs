# envs

envs is an experimental Go library that allows you to convert a standard environment variable (env) file into an encrypted envs file, which is securely stored on disk using AES-256 encryption. The library returns a hashmap of the environment variables for easy access in your applications. 

### **This should not be used in any prod.**

## Features

- Convert standard env files into encrypted envs files.
- Securely store environment variables on disk using AES-256 encryption.
- Retrieve environment variables as a hashmap for easy integration into your applications.
