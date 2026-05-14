/**
 * Transaction Processor Module
 * Handles double-spend prevention using Idempotency Keys and distributed locking patterns.
 */

export class TransactionProcessor {
    // In memory store for demonstration. In production, this must be Redis.
    private idempotencyCache: Map<string, any>;
    private accountBalances: Map<string, number>;

    constructor() {
        this.idempotencyCache = new Map();
        // Mock DB balances
        this.accountBalances = new Map();
        this.accountBalances.set("USER_A", 1000.00);
        this.accountBalances.set("USER_B", 500.00);
    }

    async processTransfer(idempotencyKey: string, fromUser: string, toUser: string, amount: number): Promise<any> {
        console.log(`[TxProcessor] Attempting transfer of $${amount} from ${fromUser} to ${toUser}`);
        
        // 1. Idempotency Check (Prevent Double Spend)
        if (this.idempotencyCache.has(idempotencyKey)) {
            console.log(`[TxProcessor] Idempotency Cache HIT for key: ${idempotencyKey}. Returning previous result.`);
            return this.idempotencyCache.get(idempotencyKey);
        }

        // 2. Validate Funds
        const senderBalance = this.accountBalances.get(fromUser) || 0;
        if (senderBalance < amount) {
            throw new Error("Insufficient Funds");
        }

        // 3. Simulated Atomic Execution (Mutex lock would go here)
        console.log(`[TxProcessor] Executing atomic state mutation...`);
        this.accountBalances.set(fromUser, senderBalance - amount);
        this.accountBalances.set(toUser, (this.accountBalances.get(toUser) || 0) + amount);

        const receipt = {
            status: "SUCCESS",
            txId: `TX_${Date.now()}`,
            timestamp: new Date().toISOString(),
            amount: amount,
            newBalance: this.accountBalances.get(fromUser)
        };

        // 4. Save to Idempotency Cache
        this.idempotencyCache.set(idempotencyKey, receipt);
        
        return receipt;
    }

    getBalance(user: string): number {
        return this.accountBalances.get(user) || 0;
    }
}
