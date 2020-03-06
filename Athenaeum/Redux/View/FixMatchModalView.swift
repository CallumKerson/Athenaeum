/**
 FixMatchModalView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct FixMatchModalView: View {
    @ObservedObject var viewModel: FixMatchModalViewModel

    @State var goodReadsID: String = ""

    init(viewModel: FixMatchModalViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        VStack {
            Text("Enter GoodReads ID or link to fix the match:")
                .padding(.top)
            TextField("GoodReads ID", text: $goodReadsID)
                .padding(.bottom)
            HStack {
                Button("Cancel") {
                    self.viewModel.cancelButtonAction()
                }
                Button("Fix Match") {
                    self.viewModel.fixMatchButton(goodReadsID: self.goodReadsID)
                }.disabled(goodReadsID.isEmpty)
            }
        }
        .padding()
    }
}

#if DEBUG
    struct FixMatchModalView_Previews: PreviewProvider {
        static var previews: some View {
            FixMatchModalView(viewModel: FixMatchModalViewModel(audiobook: sampleAudiobook,
                                                                store: sampleStore))
        }
    }
#endif
