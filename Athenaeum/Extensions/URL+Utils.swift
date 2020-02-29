/**
 URL+Utils.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import CryptoKit
import Foundation

extension URL {
    var deSandboxedPath: String {
        path
            .replacingOccurrences(of: "Library/Containers/com.umbra.Athenaeum/Data/",
                                  with: "")
    }

    func isSameIgnoringSandbox(as otherURL: URL) -> Bool {
        self.deSandboxedPath == otherURL.deSandboxedPath
    }

    var isDirectory: Bool {
        (try? resourceValues(forKeys: [.isDirectoryKey]))?.isDirectory ?? false
    }
}

extension URL {
    var sha256HashOfContents: String? {
        var hasher = SHA256()

        if let stream = InputStream(fileAtPath: path) {
            stream.open()
            let bufferSize = 2048
            let buffer = UnsafeMutablePointer<UInt8>.allocate(capacity: bufferSize)
            while stream.hasBytesAvailable {
                let read = stream.read(buffer, maxLength: bufferSize)
                let bufferPointer = UnsafeRawBufferPointer(start: buffer,
                                                           count: read)
                hasher.update(bufferPointer: bufferPointer)
            }
            return hasher.finalize().stringValue
        }
        return nil
    }
}
