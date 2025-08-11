
# Example promise sequence

```mermaid
sequenceDiagram
    participant Client
    participant Server

    Note over Client,Server: Initial State: Client has a capability for LoginService.

    %% 1. Client calls login()
    Client->>Server: 1. call(questionId: 1, target: LoginService, method: login)
    Note right of Client: Client creates a promise<br/>internally for the result<br/>of question 1.

    %% 2. Server immediately returns a promise
    Server->>Client: 2. return(answerId: 1, results: {session: promise<55>})
    Note left of Server: Server starts auth check.<br/>It creates a new promise placeholder<br/>on its side with ID 55.

    %% Server's async work
    Note over Server: ... (Database lookup, session creation) ...

    %% 3. Server resolves the promise
    Note left of Server: Auth success! A real UserSession<br/>object is created (gets exportId: 88).
    Server->>Client: 3. resolve(promiseId: 55, resolution: {cap: {id: 88}})
    Note right of Client: Client now knows that its<br/>promise for question 1 is<br/>the real remote object 88.

    %% 4. Client uses the resolved session
    Client->>Server: 4. call(questionId: 2, target: {importedCap: {id: 88}}, method: getProfile)

    %% 5. Server responds to the new call
    Server->>Client: 5. return(answerId: 2, results: {name: "Alice", ...})

    Note over Client: ... (Client is done with the session) ...

    %% 6. Client releases the session capability
    Client->>Server: 6. release(id: 88, referenceCount: 1)
    Note left of Server: Server decrements ref count for object 88.<br/>If it hits zero, the session is garbage collected.
    
```


# Example Orchestrated 3PH Sequence

```mermaid
sequenceDiagram
    participant Alice as Alice (Introducer)
    participant Bob as Bob (Receiver)
    participant Carol as Carol (Provider)

    Note over Alice, Carol: Setup: Alice has capabilities for both Bob and Carol.<br/>Bob and Carol do not have a direct connection.

    %% Step 1: Alice tells Bob to expect a capability from Carol.
    Alice->>Bob: 1. provide(questionId: 10, target: cap_to_Carol, recipientId: secret_token)
    Note right of Alice: Alice is asking Bob to<br/>initiate contact with Carol.

    Note left of Bob: Bob now knows he needs to contact<br/>Carol and present the 'secret_token'.<br/>He creates a promise for the eventual result.

    %% Step 2: Bob contacts Carol to accept the introduction.
    Bob->>Carol: 2. accept(questionId: 10, provisionId: secret_token, embargo: true)
    Note right of Bob: 'questionId' links to the 'provide'.<br/>'provisionId' proves he is the<br/>intended recipient. He sends an<br/>"embargoed" capability to himself.

    Note left of Carol: Carol verifies the 'provisionId'.<br/>She now has a (temporarily locked)<br/>way to send the result back to Bob.

    %% Step 3: Carol creates the resource and sends the capability to Bob.
    Note over Carol: ... Carol creates the DataStream resource ...
    Carol-->>Bob: 3. return(answerId: 10, results: join({joinResult: cap_to_DataStream}))
    Note right of Carol: Carol resolves the promise from step 2.<br/>The `join` message contains the<br/>final, real capability for Bob.

    %% Step 4: The final state
    Note over Bob, Carol: Handoff Complete! Bob now holds a direct<br/>capability for Carol's DataStream.
    Bob->>Carol: 4. (Direct communication using the new capability)

    Note over Alice: Alice is no longer involved in the communication<br/>between Bob and Carol.
```


# Example Promise Resolution 3PH

## The Scenario

*   **Vat A (Alice):** Our client. She wants to get a capability and immediately use it.
*   **Vat B (The Intermediary):** A service that acts as a factory or broker. It doesn't host the final capability itself.
*   **Vat C (The Provider/Carol):** The service that actually owns and hosts the final resource.

## Diagram

```mermaid
sequenceDiagram
    participant Vat_A as Vat A (Alice)
    participant Vat_B as Vat B (Intermediary)
    participant Vat_C as Vat C (Carol/Provider)

    %% PHASE 1: Setup and Pipelining
    Note over Vat_A, Vat_B: Phase 1: Alice gets a promise from B and immediately pipelines a call on it.
    Vat_A->>Vat_B: 1. call(qId:1, getCapability())
    Note right of Vat_A: "B, give me a capability."

    Vat_B->>Vat_A: 2. return(aId:1, results:{cap: senderPromise(id:55)})
    Note left of Vat_B: "OK, here is a promise for it.<br/>I'll resolve it later."

    Vat_A->>Vat_B: 3. call(qId:2, target:promisedAnswer(qId:1), method:foo())
    Note right of Vat_A: "Thanks. I'll call foo() on that<br/>promise right now."
    Note left of Vat_B: B receives the pipelined call for foo()<br/>and holds it, waiting for promise 55 to resolve.

    %% PHASE 2: Promise Resolution and Redirection
    Note over Vat_B, Vat_C: Phase 2: B realizes the capability is on C and initiates the handoff.
    Vat_B->>Vat_C: 4. provide(qId:3, for_recipient:Vat A)
    Note right of Vat_B: B tells C: "Prepare to provide a<br/>capability to Vat A. This is context 3."

    Vat_B->>Vat_A: 5. resolve(promiseId:55, resolution:ThirdPartyCapId(for C, contextId:3))
    Note left of Vat_B: B tells A: "By the way, that promise 55 is<br/>actually on Vat C. Use context 3 to get it."

    %% PHASE 3: The Handoff
    Note over Vat_A, Vat_C: Phase 3: Alice contacts Carol directly to complete the handoff.
    Vat_A->>Vat_C: 6. accept(qId:3, embargo:true)
    Note right of Vat_A: "C, I'm here for context 3. But wait!<br/>Put an embargo on the result because I<br/>have an outstanding call (foo)."
    Note left of Vat_C: "C does *not* return anything yet, because this <br/> pipeline is embargoed. API-wise, Alice has a promise to the <br/> result of this call."




    %% PHASE 4: Direct call.
    Note over Vat_A, Vat_C: Phase 4: Path-shortened call, processing embargoed on C.
    Note right of Vat_A: Alice now makes a NEW call.<br/>It goes directly to C, bypassing B entirely.
    Vat_A->>Vat_C: 8. call(qId:4, target:RealCap, method:bar())    


    %% %%%%%%%%%%%%%%%%%%%%%%%%%%%
    Note over Vat_A, Vat_C: XXXXXXXXXXXXXXXXXXXXXXXX end of story.

    Vat_C->>Vat_A: 7. return(aId:3, results:join({joinResult:RealCap}))
    Note left of Vat_C: "Verified. Here is your direct capability,<br/>RealCap. It's currently embargoed."

    %% PHASE 4 Forwarding and Disembargo in one step
    Note over Vat_B, Vat_C: Phase 4: B forwards the foo() call *inside* a Disembargo message to C.
    Vat_B->>Vat_C: 8. disembargo(context: {call: [Original foo() call from Step 3]})
    Note right of Vat_B: "C, here is a call that was<br/>pipelined to me. Please deliver it<br/>to the target it's now resolved to."
    
    Note left of Vat_C: C receives the Disembargo.<br/>1. It looks up the target of the enclosed call.<br/>2. Finds it's the embargoed RealCap from qId 3.<br/>3. Un-embargoes RealCap and executes foo().

    %% PHASE 5: Finalization and Path Shortening
    Note over Vat_A, Vat_C: Phase 5: The original call completes, and the path is now shortened.
    Vat_C->>Vat_A: 9. return(aId:2, results:{foo_result})
    Note left of Vat_C: C sends the result of foo() DIRECTLY to A.

    Vat_A->>Vat_C: 10. finish(qId:2)

    Vat_A->>Vat_C: 11. call(qId:4, target:RealCap, method:bar())
    Note right of Vat_A: Alice now makes a NEW call.<br/>It goes directly to C, bypassing B entirely.

    Vat_C->>Vat_A: 12. return(aId:4, results:{bar_result})
    Vat_A->>Vat_C: 13. finish(qId:4)

```

