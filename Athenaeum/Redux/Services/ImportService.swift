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

/// Moves the audiobook to the supplied library.
///
/// If the audiobook has no title and author, then the file is moved into the base path of the library.
///
/// If the audiobook has a title and author, then the file is moved to
/// `LIBRARY/{author}/{title}`
///
/// Additionally, if the audioobok has a series then the file is moved to
/// `LIBRARY/{author}/{seriesTitle}/{seriesEntry} {title}`
/// - Parameters:
///   - audiobook: Audiobook to move. The location is mutated to the location of the moved file
///   - libraryURL: The base URL of the library
func moveAudiobookToLibrary(_ audiobook: inout AudioBook, libraryURL: URL) throws {
    var destination = libraryURL
        .appendingPathComponent(audiobook.location.lastPathComponent)
    if let metadata = audiobook.metadata {
        if let author = metadata.authors?.author.replacingOccurrences(of: ":", with: " -") {
            let title = metadata.title.replacingOccurrences(of: ":", with: " -")
            destination = libraryURL
                .appendingPathComponent(author, isDirectory: true)
            if let series = metadata.series {
                destination = destination
                    .appendingPathComponent(series.title, isDirectory: true)
                    .appendingPathComponent("\(series.entry.asString) \(title)", isDirectory: false)
                    .appendingPathExtension("m4b")
            } else {
                destination = destination
                    .appendingPathComponent(title, isDirectory: false)
                    .appendingPathExtension("m4b")
            }
        }
    }

    if audiobook.location.isSameIgnoringSandbox(as: destination) {
        log.debug("New Audiobook \(audiobook) is already in the right place")
    } else {
        log.debug("Moving audiobook file to \(destination.path)")
        try FileManager.default
            .moveItemCreatingIntermediaryDirectories(at: audiobook.location, to: destination)
        audiobook.location = destination
    }
}
