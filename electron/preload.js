import { contextBridge, ipcRenderer } from "electron";

contextBridge.exposeInMainWorld("api", {
  openCsv: () => ipcRenderer.invoke("csv:open"),
});
