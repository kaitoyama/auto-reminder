# コマンドハンドラのドキュメント

このドキュメントは、traQボットのために `internal/interface/handler.go` に実装されたコマンドハンドリングロジックについて説明します。

## 概要

ハンドラはボットに送信されたメッセージを解析し、特定のキーワードとオプションに基づいてコマンドを実行します。各コマンドは `@BotName`（*BotName* は設定から取得）で始まり、続いてコマンド名とそのパラメータが続きます。コマンドの形式が不正または必要なオプションが欠けている場合、ボットはエラーメッセージ「コマンドが不正です」を返します。

## サポートされているコマンド

### create

- **使用法:**  
  `@BotName create -u @user1 @user2 ... -d YYYY-MM-DD -c コマンド内容`  
- **パラメータ:**
  - `-u`: 1人以上のユーザーをメンション（例：`@user1`）。ユーザー名は `@` の後のテキストから抽出されます。
  - `-d`: `YYYY-MM-DD` 形式の締め切り日。この日付は日付として解析されます。
  - `-c`: タスクの内容。`-c` の後のすべてのトークンが結合されてコマンド内容を形成します。
- **処理:**  
  `create.Create` ユースケースを使用してコマンドを処理します。パラメータが欠けているか無効な場合、ボットはエラーメッセージを投稿します。

### addUser

- **使用法:**  
  `@BotName addUser -u @user1 @user2 ... -i taskID`  
- **パラメータ:**
  - `-u`: 1人以上のユーザーをメンション。
  - `-i`: タスクを表す整数ID。
- **処理:**  
  `create.AddUser` ユースケースを使用します。パラメータが欠けているか無効な場合、エラーメッセージが表示されます。

### deleteUser

- **使用法:**  
  `@BotName deleteUser -u @user -i taskID`  
- **パラメータ:**
  - `-u`: 1人のユーザーをメンション（1人のみ許可）。
  - `-i`: タスクの整数ID。
- **処理:**  
  `create.DeleteUser` ユースケースを使用して指定されたユーザーを削除します。このコマンドには正確に1人のユーザーと有効なタスクIDが必要です。

### updateDueAt

- **使用法:**  
  `@BotName updateDueAt -d YYYY-MM-DD -i taskID`  
- **パラメータ:**
  - `-d`: `YYYY-MM-DD` 形式の新しい締め切り日。
  - `-i`: タスクの整数ID。
- **処理:**  
  `create.UpdateDueAt` ユースケースを使用してタスクの締め切りを更新します。

### delete

- **使用法:**  
  `@BotName delete -i taskID`  
- **パラメータ:**
  - `-i`: タスクの整数ID。
- **処理:**  
  `create.Delete` ユースケースを使用して指定されたタスクを削除します。

## エラーハンドリング

すべてのコマンドについて、次の場合にボットはエラーメッセージ「コマンドが不正です」を投稿します:
- 必須オプションやパラメータが欠けている場合、
- 提供されたオプションが解析できない場合（例：日付形式が不正または整数IDでない）、
- コマンド名がサポートされているコマンドと一致しない場合。

## フローの概要

1. **メッセージ受信:**  
   ボットは最初に受信したメッセージが `@BotName` で始まることを確認します。

2. **引数解析:**  
   メッセージテキストはトークンに分割されます。コマンドと各オプションは特定のフラグに基づいて識別されます。

3. **検証とユースケースの実行:**  
   - 有効なコマンドとオプションの場合、対応するユースケースが実行されます。
   - 解析中またはユースケースの実行中にエラーが発生した場合、ユーザーにエラーメッセージが表示されます。

4. **ログ記録:**  
   メッセージ受信やエラーなどの重要なイベントはすべて [zerolog](https://github.com/rs/zerolog) を使用してログに記録されます。

このドキュメントは、コマンドハンドラのコア機能とエラーハンドリング戦略をカバーしています。

# Command Handler Documentation

This document explains the command handling logic implemented in `internal/interface/handler.go` for the traQ bot.

## Overview

The handler parses messages sent to the bot and executes commands based on specific keywords and options. Each command is expected to be prefixed with `@BotName` (where *BotName* is obtained from configuration) followed by the command name and its parameters. If the command format is incorrect or required options are missing, the bot responds with an error message: 「コマンドが不正です」.

## Supported Commands

### create

- **Usage:**  
  `@BotName create -u @user1 @user2 ... -d YYYY-MM-DD -c command content`  
- **Parameters:**
  - `-u`: One or more user mentions (e.g., `@user1`). The username is extracted (text after `@`).
  - `-d`: Deadline in the format `YYYY-MM-DD`. This is parsed into a date.
  - `-c`: Content for the task. All remaining tokens after `-c` are joined to form the command content.
- **Processing:**  
  Uses the `create.Create` usecase to process the command. If any parameter is missing or invalid, the bot posts an error message.

### addUser

- **Usage:**  
  `@BotName addUser -u @user1 @user2 ... -i taskID`  
- **Parameters:**
  - `-u`: One or more user mentions.
  - `-i`: An integer ID representing the task.
- **Processing:**  
  Uses the `create.AddUser` usecase. Missing or invalid parameters result in an error message.

### deleteUser

- **Usage:**  
  `@BotName deleteUser -u @user -i taskID`  
- **Parameters:**
  - `-u`: A single user mention (only one user allowed).
  - `-i`: An integer ID for the task.
- **Processing:**  
  Uses the `create.DeleteUser` usecase to remove the specified user. The command requires exactly one user and a valid task ID.

### updateDueAt

- **Usage:**  
  `@BotName updateDueAt -d YYYY-MM-DD -i taskID`  
- **Parameters:**
  - `-d`: New deadline date in the format `YYYY-MM-DD`.
  - `-i`: An integer ID for the task.
- **Processing:**  
  Uses the `create.UpdateDueAt` usecase to update the task's deadline.

### delete

- **Usage:**  
  `@BotName delete -i taskID`  
- **Parameters:**
  - `-i`: An integer ID for the task.
- **Processing:**  
  Uses the `create.Delete` usecase to remove the specified task.

## Error Handling

For all commands, if:
- Mandatory options or parameters are missing,
- The provided options cannot be parsed (e.g., incorrect date format or non-integer ID), or
- The command name does not match any supported commands,

the bot posts the error message: 「コマンドが不正です」 indicating that the command is invalid.

## Flow Summary

1. **Message Reception:**  
   The bot first verifies that the received message starts with `@BotName`.

2. **Argument Parsing:**  
   The message text is split into tokens. The command and each option are identified based on specific flags.

3. **Validation and Usecase Execution:**  
   - For a valid command and options, the corresponding usecase is executed.
   - Errors during parsing or usecase execution result in an error message to the user.

4. **Logging:**  
   All significant events, such as message receipt and errors, are logged using [zerolog](https://github.com/rs/zerolog).

This documentation covers the core functionality and error handling strategy of the command handler.
