/**
 Librarian+Cocoa.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Cocoa
import Foundation

extension Librarian {
    func openImportAudiobookDialog() {
        let openPanel = NSOpenPanel()
        openPanel.canChooseFiles = true
        openPanel.allowsMultipleSelection = true
        openPanel.canChooseDirectories = false
        openPanel.canCreateDirectories = false
        openPanel.title = "Import Audiobooks"
        openPanel.allowedFileTypes = ["m4b"]

        openPanel.begin { response in
            if response == .OK {
                for url in openPanel.urls {
                    DispatchQueue.global(qos: .userInitiated).async {
                        log.debug("User opening file \(url.path)")
                        self.importAudiobook(fileURL: url)
                    }
                }
            }
            openPanel.close()
        }
    }
}
