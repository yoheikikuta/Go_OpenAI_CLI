# OpenAI Chatbot

This repository contains a simple chatbot that interacts with the OpenAI GPT-3.5-turbo model in real-time. The chatbot takes user input and sends it to the OpenAI API, which returns a generated response that the chatbot displays in the terminal.

<div align="center">
  <video src="https://user-images.githubusercontent.com/14804123/231175566-2dbcef22-2233-4b30-a70d-f1d2747a9b25.mp4" width="400" />
</div>

## Features
Real-time interaction with the OpenAI API
Simple command-line interface
Easy to understand code structure, suitable for customization

## Requirements 
Go 1.16 or later  
An OpenAI API key

## Setup
Clone this repository:

```bash
git clone https://github.com/yoheikikuta/Go_OpenAI_CLI.git
cd Go_OpenAI_CLI
```

Create a file named api_key.txt in the project directory and paste your OpenAI API key inside.


Build and run the chatbot:

```bash
go build -o chat
./chat
```

## Usage
1. Enter your message in the terminal and press Enter.
2. The chatbot will display the generated response from the OpenAI API.
3. To exit the chatbot, type "exit" and press Enter.


## Customization
You can customize the chatbot by modifying the source code. Some customization options include:

- Adjusting the model's temperature and top_p parameters for different response characteristics
- Implementing a custom error handling mechanism
- Adding more features, like logging or user authentication

## License
This project is licensed under the MIT License. See the LICENSE file for details.
