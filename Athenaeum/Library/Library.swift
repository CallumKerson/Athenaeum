/**
 Library.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

protocol Library: ObservableObject {
    var 🎧📚: [Audiobook] { get set }
    func shelve(book: Audiobook)
}
