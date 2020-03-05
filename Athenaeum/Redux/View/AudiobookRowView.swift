/**
 AudiobookRowView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct AudiobookRowView: View {
    @ObservedObject var viewModel: AudiobookRowViewModel

    init(_ viewModel: AudiobookRowViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        Unwrap(viewModel.audiobook) { audiobook in
            HStack(alignment: VerticalAlignment.center) {
                Unwrap(audiobook.getCover()) { coverData in
                    RowCover(data: coverData)
                }

                VStack(alignment: HorizontalAlignment.leading) {
                    Unwrap(audiobook.title) { title in
                        Text(title)
                            .fontWeight(.bold)
                            .truncationMode(.tail)
                            .frame(minWidth: 20)
                    }

                    Unwrap(audiobook.getAuthorsString()) { author in
                        Text(author)
                            .font(.caption)
                            .opacity(0.625)
                            .truncationMode(.middle)
                    }
                }

                Spacer()
            }
            .padding(.vertical, 4)
        }
    }
}

struct RowCover: View {
    let data: Data

    var body: some View {
        let image = NSImage(data: data)!
        return Image(nsImage: image.resize(withSize: NSMakeSize(32, 32))!)
            .frame(width: 32, height: 32)
            .fixedSize(horizontal: true, vertical: false)
            .clipShape(RoundedRectangle(cornerRadius: 5))
            .shadow(radius: 2)
    }
}

#if DEBUG
    struct AudiobookRowView_Previews: PreviewProvider {
        static var previews: some View {
            AudiobookRowView(AudiobookRowViewModel(id: sampleAudiobook.id,
                                                   store: sampleStore))
        }
    }
#endif
