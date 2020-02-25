/**
 Audiobook.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

protocol Audiobook {
    var title: String { get set }
    var author: String { get set }
    var location: URL { get set }
    var narrator: String? { get set }
    var publicationDate: String? { get set }
    var isbn: String? { get set }
    var summary: String? { get set }
    var series: Series? { get set }

    func getCover() -> Data?
}
