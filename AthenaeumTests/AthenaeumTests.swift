/**
 AthenaeumTests.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

@testable import Athenaeum
import AVFoundation
import XCTest

class AthenaeumTests: XCTestCase {
    override func setUp() {
        // Put setup code here. This method is called before the invocation of each test method in the class.
    }

    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
    }

    func testExample() {
//        let audiobook = Audiobook(fileURL: URL(string: "/Users/ckerson/Music/TWoK.m4b")!)

        let newBook = Audiobook.getBookFromFile(path: "/Users/ckerson/Music/TWoK.m4b")
        let destination = Library.global.libraryURL
            .appendingPathComponent(newBook.author, isDirectory: true)
            .appendingPathComponent(newBook.title)
            .appendingPathExtension("m4b")
        print(destination)
    }

    func testPerformanceExample() {
        // This is an example of a performance test case.
        measure {
            // Put the code you want to measure the time of here.
        }
    }
}
