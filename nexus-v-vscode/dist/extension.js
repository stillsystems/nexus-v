"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = require("vscode");
const path = require("path");
const fs = require("fs");
function activate(context) {
    const disposable = vscode.commands.registerCommand('nexus-v.createProject', () => {
        const panel = vscode.window.createWebviewPanel('nexusVScaffolder', 'Nexus-V: New Project', vscode.ViewColumn.One, {
            enableScripts: true,
            localResourceRoots: [vscode.Uri.file(path.join(context.extensionPath, 'media'))]
        });
        panel.webview.html = getWebviewContent(context, panel.webview);
        // Handle messages from the webview
        panel.webview.onDidReceiveMessage((message) => {
            switch (message.command) {
                case 'generate':
                    vscode.window.showInformationMessage(`Generating project: ${message.data.name}...`);
                    // In a production version, we would spawn the nexus-v binary here
                    return;
            }
        }, undefined, context.subscriptions);
    });
    context.subscriptions.push(disposable);
}
function getWebviewContent(context, webview) {
    const htmlPath = path.join(context.extensionPath, 'media', 'index.html');
    let html = fs.readFileSync(htmlPath, 'utf8');
    // Get URIs for local resources
    const styleUri = webview.asWebviewUri(vscode.Uri.file(path.join(context.extensionPath, 'media', 'style.css')));
    const scriptUri = webview.asWebviewUri(vscode.Uri.file(path.join(context.extensionPath, 'media', 'app.js')));
    // Replace relative paths with Webview URIs
    html = html.replace('href="style.css"', `href="${styleUri}"`);
    html = html.replace('src="app.js"', `src="${scriptUri}"`);
    // Add VS Code API acquire snippet
    html = html.replace('</head>', `
        <script>
            const vscode = acquireVsCodeApi();
            window.vscode = vscode;
        </script>
    </head>`);
    return html;
}
function deactivate() { }
//# sourceMappingURL=extension.js.map