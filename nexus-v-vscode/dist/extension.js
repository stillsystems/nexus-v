"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = require("vscode");
const path = require("path");
const fs = require("fs");
function activate(context) {
    const controlProvider = new NexusControlProvider();
    vscode.window.registerTreeDataProvider('nexusv.dashboard', controlProvider);
    context.subscriptions.push(vscode.commands.registerCommand('nexusv.refresh', () => controlProvider.refresh()));
    context.subscriptions.push(vscode.commands.registerCommand('nexus-v.runDoctor', () => {
        const terminal = vscode.window.createTerminal("Nexus-V Doctor");
        terminal.show();
        terminal.sendText("nexus-v doctor");
    }));
    const disposable = vscode.commands.registerCommand('nexus-v.createProject', () => {
        // ... (existing code remains)
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
class NexusControlProvider {
    constructor() {
        this._onDidChangeTreeData = new vscode.EventEmitter();
        this.onDidChangeTreeData = this._onDidChangeTreeData.event;
    }
    refresh() {
        this._onDidChangeTreeData.fire();
    }
    getTreeItem(element) {
        return element;
    }
    getChildren(element) {
        if (element) {
            return Promise.resolve(element.children || []);
        }
        else {
            return Promise.resolve([
                new NexusItem('Project Status', vscode.TreeItemCollapsibleState.Expanded, 'dashboard', [
                    new NexusItem('Engine: v0.2.8', vscode.TreeItemCollapsibleState.None, 'check'),
                    new NexusItem('Health: Perfect', vscode.TreeItemCollapsibleState.None, 'shield')
                ]),
                new NexusItem('Quick Actions', vscode.TreeItemCollapsibleState.Expanded, 'zap', [
                    new NexusItem('Create New Project', vscode.TreeItemCollapsibleState.None, 'add', undefined, {
                        command: 'nexus-v.createProject',
                        title: 'New Project'
                    }),
                    new NexusItem('Run Doctor', vscode.TreeItemCollapsibleState.None, 'pulse', undefined, {
                        command: 'nexus-v.runDoctor',
                        title: 'Run Doctor'
                    })
                ])
            ]);
        }
    }
}
class NexusItem extends vscode.TreeItem {
    constructor(label, collapsibleState, iconName, children, command) {
        super(label, collapsibleState);
        this.label = label;
        this.collapsibleState = collapsibleState;
        this.iconName = iconName;
        this.children = children;
        this.command = command;
        this.contextValue = 'nexusItem';
        if (iconName) {
            this.iconPath = new vscode.ThemeIcon(iconName);
        }
    }
}
//# sourceMappingURL=extension.js.map