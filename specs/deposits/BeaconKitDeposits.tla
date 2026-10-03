--------------------- MODULE BeaconKitDeposits ----------------------
EXTENDS Naturals, Sequences, TLC
CONSTANTS MAX_DEPOSITS_PER_BLOCK, N, N_DEPOSITS
Min(x, y) == IF x < y THEN x ELSE y
Tern(x, y, z) == IF x = TRUE THEN y ELSE z 
Range(f) == {f[x] : x \in DOMAIN f}
(*
--algorithm BeaconKitDeposits {
    variables
        Eth1Blocks = <<1>>;
        DepositLogs = [x \in 1..N |-> <<>>];
        DepositStore = <<>>;
        BeaconBlocks = <<>>;
        BeaconBlock = [blockNum |-> 0, eth1BlockNum |-> 0, deposits |-> <<>>];
        Finalized = FALSE;
        DepositTxs = <<>>;

    define {
        I_EmptyStore == DepositStore = <<>> <=> \A ds \in Range(DepositLogs) : ds = <<>>
        P_EmptyStore == []<>I_EmptyStore
        I == \A ds \in Range(SubSeq(DepositLogs, 1, Len(DepositLogs)-1)) : ds /= <<>> => \E d \in Range(DepositStore) : d \in Range(ds)
        P == []<>I
    }

    fair process (DepositContract = "deposit-contract")
        variables
            depositId = 0;
    {
        deposited:
            while (depositId < N_DEPOSITS) {
                either {
                    depositId := depositId + 1;
                    DepositTxs := Append(DepositTxs, depositId);   \* deposit contract received tx
                } or {
                    skip;
                }
            }
    }

    fair process (ExecutionLayer = "execution-layer")
        variables
            eth1BlockNumber = 1;
    {
        produceEth1Block:
            while (eth1BlockNumber < N) {
                either {
                    eth1BlockNumber := eth1BlockNumber + 1;
                    DepositLogs[eth1BlockNumber] := DepositTxs;
                    Eth1Blocks := Append(Eth1Blocks, eth1BlockNumber);
                    DepositTxs := <<>>;
                } or {
                    skip;
                };
            };
        \* end produceEth1Block
    }

    fair process (BeaconBlockProducer = "beacon-block-producer")
        variables
            beaconBlock = BeaconBlock;
            deposits = <<>>;
            depositIndex = 0;
            numDeposits = 0;
            eth1N = 0;
    {
        emitFinalized:
            while (TRUE) {
                await (Finalized = FALSE /\ Eth1Blocks /= <<>>);
                eth1N := Head(Eth1Blocks);
                Eth1Blocks := Tail(Eth1Blocks);
                numDeposits := Min(Len(DepositStore) - depositIndex, MAX_DEPOSITS_PER_BLOCK);
                deposits := Tern(depositIndex > 0, SubSeq(DepositStore, depositIndex, numDeposits), <<>>);
                beaconBlock := [beaconBlock EXCEPT !.blockNum = @ + 1, !.eth1BlockNum = eth1N, !.deposits = deposits];
                BeaconBlocks := Append(BeaconBlocks, beaconBlock);
                depositIndex := depositIndex + numDeposits;
                Finalized := TRUE;
            }
        \* end emitFinalized
    }

    fair process (DepositFetcher = "deposit-fetcher")
        variables
            pendingDeposits = <<>>;
            curBlock = BeaconBlock;
    {
        fetch:
            while (TRUE) {
                await (Finalized = TRUE);
                curBlock := BeaconBlocks[Len(BeaconBlocks)];
                if (curBlock.eth1BlockNum > 1) {
                    pendingDeposits := DepositLogs[curBlock.eth1BlockNum-1];
                    DepositStore := DepositStore \o pendingDeposits;
                };
                Finalized := FALSE;
            };
    }

} \* end algorithm
*)
\* BEGIN TRANSLATION (chksum(pcal) = "62c5d1d9" /\ chksum(tla) = "2befdb2f")
VARIABLES pc, Eth1Blocks, DepositLogs, DepositStore, BeaconBlocks, 
          BeaconBlock, Finalized, DepositTxs

(* define statement *)
I_EmptyStore == DepositStore = <<>> <=> \A ds \in Range(DepositLogs) : ds = <<>>
P_EmptyStore == []<>I_EmptyStore
I == \A ds \in Range(SubSeq(DepositLogs, 1, Len(DepositLogs)-1)) : ds /= <<>> => \E d \in Range(DepositStore) : d \in Range(ds)
P == []<>I

VARIABLES depositId, eth1BlockNumber, beaconBlock, deposits, depositIndex, 
          numDeposits, eth1N, pendingDeposits, curBlock

vars == << pc, Eth1Blocks, DepositLogs, DepositStore, BeaconBlocks, 
           BeaconBlock, Finalized, DepositTxs, depositId, eth1BlockNumber, 
           beaconBlock, deposits, depositIndex, numDeposits, eth1N, 
           pendingDeposits, curBlock >>

ProcSet == {"deposit-contract"} \cup {"execution-layer"} \cup {"beacon-block-producer"} \cup {"deposit-fetcher"}

Init == (* Global variables *)
        /\ Eth1Blocks = <<1>>
        /\ DepositLogs = [x \in 1..N |-> <<>>]
        /\ DepositStore = <<>>
        /\ BeaconBlocks = <<>>
        /\ BeaconBlock = [blockNum |-> 0, eth1BlockNum |-> 0, deposits |-> <<>>]
        /\ Finalized = FALSE
        /\ DepositTxs = <<>>
        (* Process DepositContract *)
        /\ depositId = 0
        (* Process ExecutionLayer *)
        /\ eth1BlockNumber = 1
        (* Process BeaconBlockProducer *)
        /\ beaconBlock = BeaconBlock
        /\ deposits = <<>>
        /\ depositIndex = 0
        /\ numDeposits = 0
        /\ eth1N = 0
        (* Process DepositFetcher *)
        /\ pendingDeposits = <<>>
        /\ curBlock = BeaconBlock
        /\ pc = [self \in ProcSet |-> CASE self = "deposit-contract" -> "deposited"
                                        [] self = "execution-layer" -> "produceEth1Block"
                                        [] self = "beacon-block-producer" -> "emitFinalized"
                                        [] self = "deposit-fetcher" -> "fetch"]

deposited == /\ pc["deposit-contract"] = "deposited"
             /\ IF depositId < N_DEPOSITS
                   THEN /\ \/ /\ depositId' = depositId + 1
                              /\ DepositTxs' = Append(DepositTxs, depositId')
                           \/ /\ TRUE
                              /\ UNCHANGED <<DepositTxs, depositId>>
                        /\ pc' = [pc EXCEPT !["deposit-contract"] = "deposited"]
                   ELSE /\ pc' = [pc EXCEPT !["deposit-contract"] = "Done"]
                        /\ UNCHANGED << DepositTxs, depositId >>
             /\ UNCHANGED << Eth1Blocks, DepositLogs, DepositStore, 
                             BeaconBlocks, BeaconBlock, Finalized, 
                             eth1BlockNumber, beaconBlock, deposits, 
                             depositIndex, numDeposits, eth1N, pendingDeposits, 
                             curBlock >>

DepositContract == deposited

produceEth1Block == /\ pc["execution-layer"] = "produceEth1Block"
                    /\ IF eth1BlockNumber < N
                          THEN /\ \/ /\ eth1BlockNumber' = eth1BlockNumber + 1
                                     /\ DepositLogs' = [DepositLogs EXCEPT ![eth1BlockNumber'] = DepositTxs]
                                     /\ Eth1Blocks' = Append(Eth1Blocks, eth1BlockNumber')
                                     /\ DepositTxs' = <<>>
                                  \/ /\ TRUE
                                     /\ UNCHANGED <<Eth1Blocks, DepositLogs, DepositTxs, eth1BlockNumber>>
                               /\ pc' = [pc EXCEPT !["execution-layer"] = "produceEth1Block"]
                          ELSE /\ pc' = [pc EXCEPT !["execution-layer"] = "Done"]
                               /\ UNCHANGED << Eth1Blocks, DepositLogs, 
                                               DepositTxs, eth1BlockNumber >>
                    /\ UNCHANGED << DepositStore, BeaconBlocks, BeaconBlock, 
                                    Finalized, depositId, beaconBlock, 
                                    deposits, depositIndex, numDeposits, eth1N, 
                                    pendingDeposits, curBlock >>

ExecutionLayer == produceEth1Block

emitFinalized == /\ pc["beacon-block-producer"] = "emitFinalized"
                 /\ (Finalized = FALSE /\ Eth1Blocks /= <<>>)
                 /\ eth1N' = Head(Eth1Blocks)
                 /\ Eth1Blocks' = Tail(Eth1Blocks)
                 /\ numDeposits' = Min(Len(DepositStore) - depositIndex, MAX_DEPOSITS_PER_BLOCK)
                 /\ deposits' = Tern(depositIndex > 0, SubSeq(DepositStore, depositIndex, numDeposits'), <<>>)
                 /\ beaconBlock' = [beaconBlock EXCEPT !.blockNum = @ + 1, !.eth1BlockNum = eth1N', !.deposits = deposits']
                 /\ BeaconBlocks' = Append(BeaconBlocks, beaconBlock')
                 /\ depositIndex' = depositIndex + numDeposits'
                 /\ Finalized' = TRUE
                 /\ pc' = [pc EXCEPT !["beacon-block-producer"] = "emitFinalized"]
                 /\ UNCHANGED << DepositLogs, DepositStore, BeaconBlock, 
                                 DepositTxs, depositId, eth1BlockNumber, 
                                 pendingDeposits, curBlock >>

BeaconBlockProducer == emitFinalized

fetch == /\ pc["deposit-fetcher"] = "fetch"
         /\ (Finalized = TRUE)
         /\ curBlock' = BeaconBlocks[Len(BeaconBlocks)]
         /\ IF curBlock'.eth1BlockNum > 1
               THEN /\ pendingDeposits' = DepositLogs[curBlock'.eth1BlockNum-1]
                    /\ DepositStore' = DepositStore \o pendingDeposits'
               ELSE /\ TRUE
                    /\ UNCHANGED << DepositStore, pendingDeposits >>
         /\ Finalized' = FALSE
         /\ pc' = [pc EXCEPT !["deposit-fetcher"] = "fetch"]
         /\ UNCHANGED << Eth1Blocks, DepositLogs, BeaconBlocks, BeaconBlock, 
                         DepositTxs, depositId, eth1BlockNumber, beaconBlock, 
                         deposits, depositIndex, numDeposits, eth1N >>

DepositFetcher == fetch

Next == DepositContract \/ ExecutionLayer \/ BeaconBlockProducer
           \/ DepositFetcher

Spec == /\ Init /\ [][Next]_vars
        /\ WF_vars(DepositContract)
        /\ WF_vars(ExecutionLayer)
        /\ WF_vars(BeaconBlockProducer)
        /\ WF_vars(DepositFetcher)

\* END TRANSLATION 
==================================================================
