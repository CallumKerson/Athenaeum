/**
 AudiobookListView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation
import SwiftUI

struct AudiobookListView: View {
    @ObservedObject var viewModel: AudiobookListViewModel

    init(_ viewModel: AudiobookListViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        List(selection: $viewModel.selectedAudiobook) {
            ForEach(viewModel.audiobooks.sorted(by: {
                if let authorZero = $0.authors?.first, let authorOne = $1.authors?.first {
                    if authorZero != authorOne {
                        return authorZero.lastName < authorOne.lastName
                    }
                }
                if let pubDateZero = $0.publicationDate,
                    let pubDateOne = $1.publicationDate {
                    return pubDateZero < pubDateOne
                } else {
                    return $0.title < $1.title
                }
            }), id: \.id) { audiobook in
                AudiobookRowView(AudiobookRowViewModel(id: audiobook.id,
                                                       store: store))
                    .tag(audiobook)
            }
        }
    }
}

extension Binding {
    /// Wrapper to listen to didSet of Binding
    func didSet(_ didSet: @escaping ((newValue: Value, oldValue: Value)) -> Void)
        -> Binding<Value> {
        .init(get: { self.wrappedValue }, set: { newValue in
            let oldValue = self.wrappedValue
            self.wrappedValue = newValue
            didSet((newValue, oldValue))
        })
    }
}
