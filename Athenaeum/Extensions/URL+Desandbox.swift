/**
 URL+Desandbox.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

extension URL {
    func deSandboxedPath() -> String {
        path.replacingOccurrences(of: "Library/Containers/com.umbra.Athenaeum/Data/", with: "")
    }
}
