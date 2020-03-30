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
                if let metadataZero = $0.metadata, let metadataOne = $1.metadata {
                    if let authorZero = metadataZero.authors?.first,
                        let authorOne = metadataOne.authors?.first {
                        if authorZero != authorOne {
                            return authorZero.lastName < authorOne.lastName
                        }
                    }
                    if let pubDateZero = metadataZero.publicationDate?.asDate,
                        let pubDateOne = metadataOne.publicationDate?.asDate {
                        return pubDateZero < pubDateOne
                    }
                    return metadataZero.title < metadataOne.title
                } else {
                    return $0.id.uuidString < $1.id.uuidString
                }
            }), id: \.id) { audiobook in
                AudiobookRowView(AudiobookRowViewModel(id: audiobook.id,
                                                       store: self.viewModel.store))
                    .tag(audiobook)
            }
        }
    }
}

#if DEBUG
    struct AudiobookListView_Previews: PreviewProvider {
        static var previews: some View {
            AudiobookListView(AudiobookListViewModel(store: sampleStore))
        }
    }
#endif
