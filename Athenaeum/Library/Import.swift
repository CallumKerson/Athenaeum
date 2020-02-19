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
    openPanel.title = "Import Audiobooks"
    openPanel.allowedFileTypes = ["m4b"]

    openPanel.begin { response in
        if response == .OK {
            for url in openPanel.urls {
                log.debug("User opening file \(url.path)")
                DispatchQueue.main.async {
                    importAudiobook(fileURL: url)
                }
            }
        }
        openPanel.close()
    }
}

func importAudiobook(fileURL: URL) {
    log.info("Importing audiobook file from \(fileURL.path)")
    let newBook = Audiobook(fromFile: fileURL)
    var destination = Preferences.libraryPath()
        .appendingPathComponent(newBook.author, isDirectory: true)
    if let series = newBook.series {
        destination = destination
            .appendingPathComponent(series.title, isDirectory: true)
            .appendingPathComponent("\(series.entry) \(newBook.title)")
            .appendingPathExtension("m4b")
    } else {
        destination = destination
            .appendingPathComponent(newBook.title)
            .appendingPathExtension("m4b")
    }
    log.debug("Moving audiobook file to \(destination.path)")
    try! FileManager.default.moveItemCreatingIntermediaryDirectories(at: fileURL, to: destination)
    newBook.file = destination
}
