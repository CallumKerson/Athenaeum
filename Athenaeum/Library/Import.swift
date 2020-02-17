/**
 Import.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Cocoa
import Foundation

func openImportAudiobookDialog() {
    let openPanel = NSOpenPanel()
    openPanel.canChooseFiles = true
    openPanel.allowsMultipleSelection = true
    openPanel.canChooseDirectories = false
    openPanel.canCreateDirectories = false
    openPanel.title = "Import Audiobook"
    openPanel.allowedFileTypes = ["m4b"]

    openPanel.begin { response in
        if response == .OK {
            for url in openPanel.urls {
                print(url.path)
                DispatchQueue.main.async {
                    Library.global.importAudiobook(fileURL: url)
                }
            }
        }
        openPanel.close()
    }
}
