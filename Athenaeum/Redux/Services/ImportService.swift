/**
 ImportService.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Cocoa
import Foundation

/// Opens the open file dialog and imports selected m4b files
/// - Parameter store: store to which to send import request
func importFromOpenDialog(store: Store<GlobalAppState>) {
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
                store
                    .dispatch(action: AudiobookActions
                        .RequestNewAudiobookFromFile(fileURL: url))
            }
        }
        openPanel.close()
    }
}
