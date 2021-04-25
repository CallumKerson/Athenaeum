/**
 AthenaeumError.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Foundation

enum AthenaeumError: Error {
    case message(String)
    case other(Error)

    static func map(_ error: Error) -> AthenaeumError {
        (error as? AthenaeumError) ?? .other(error)
    }
}

extension AthenaeumError: CustomStringConvertible {
    var description: String {
        switch self {
        case let .message(message):
            return message
        case let .other(error):
            return error.localizedDescription
        }
    }
}
