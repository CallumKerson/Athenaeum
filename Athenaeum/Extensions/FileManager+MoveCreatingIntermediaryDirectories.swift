/**
 FileManager+MoveCreatingIntermediaryDirectories.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

extension FileManager {
    func moveItemCreatingIntermediaryDirectories(at: URL, to: URL) throws {
        let parentPath = to.deletingLastPathComponent()
        if !fileExists(atPath: parentPath.path) {
            try createDirectory(at: parentPath,
                                withIntermediateDirectories: true,
                                attributes: nil)
        }
        try moveItem(at: at, to: to)
    }

    func createDirectoryIfDoesNotExist(atURL url: URL) throws {
        var isDir = ObjCBool(true)
        if fileExists(atPath: url.path, isDirectory: &isDir) == false {
            try createDirectory(at: url, withIntermediateDirectories: true,
                                attributes: nil)
        }
    }
}
