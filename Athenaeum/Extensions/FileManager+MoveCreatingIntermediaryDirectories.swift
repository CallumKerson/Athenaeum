/**
 FileManager+MoveCreatingIntermediaryDirectories.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

extension FileManager {
    func moveItemCreatingIntermediaryDirectories(at: URL, to: URL) throws {
        let parentPath = to.deletingLastPathComponent()
        if !fileExists(atPath: parentPath.path) {
            try createDirectory(at: parentPath, withIntermediateDirectories: true, attributes: nil)
        }
        try moveItem(at: at, to: to)
    }
}
