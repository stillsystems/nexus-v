import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
    console.log('Extension "Smoke Test" is now active!');

    const disposable = vscode.commands.registerCommand('smoke-test.helloWorld', () => {
        vscode.window.showInformationMessage('Smoke Test: Command executed!');
    });

    context.subscriptions.push(disposable);
}

export function deactivate() {}
