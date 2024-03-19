package tracker

// Gateway ABI

// InboundMessageDispatchedTopic InboundMessageDispatched (index_topic_1 bytes32 channelID, uint64 nonce, index_topic_2 bytes32 messageID, bool success)
const InboundMessageDispatchedTopic = "0x617fdb0cb78f01551a192a3673208ec5eb09f20a90acf673c63a0dcb11745a7a"

// OutboundMessageAcceptedTopic OutboundMessageAccepted (index_topic_1 bytes32 channelID, uint64 nonce, index_topic_2 bytes32 messageID, bytes payload)
const OutboundMessageAcceptedTopic = "0x7153f9357c8ea496bba60bf82e67143e27b64462b49041f8e689e1b05728f84f"
