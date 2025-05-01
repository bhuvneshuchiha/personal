import pandas as pd

def filter_transactions(file_path):
    # Load the Excel file
    df = pd.read_excel(file_path)

    # Drop completely empty rows
    df.dropna(how='all', inplace=True)

    # Reset index and detect header row properly
    df.columns = df.iloc[0]  # Set first non-empty row as header
    df = df.iloc[1:].reset_index(drop=True)  # Drop the old header row

    # Drop unnamed columns
    df = df.loc[:, ~df.columns.str.contains('^Unnamed', case=False)]

    # Rename columns to lowercase and strip spaces
    df.columns = [col.lower().strip() for col in df.columns]

    # Identify the narration column
    narration_cols = ['narration', 'description', 'transaction details']
    narration_col = next((col for col in narration_cols if col in df.columns), None)

    if not narration_col:
        print("Error: No valid narration column found.")
        return None

    # Convert narration column to lowercase for case-insensitive matching
    df[narration_col] = df[narration_col].astype(str).str.lower()

    # Identify amount column (could be 'amount' or 'lodgement')
    amount_cols = ['amount', 'lodgment']
    amount_col = next((col for col in amount_cols if col in df.columns), None)

    if not amount_col:
        print("Warning: No valid amount column found. Pending sheet might be incomplete.")

    # Convert amount column to numeric (if exists)
    if amount_col:
        df[amount_col] = pd.to_numeric(df[amount_col], errors='coerce')

    # Identify withdrawals column
    withdrawals_col = "withdrawals" if "withdrawals" in df.columns else None

    if withdrawals_col:
        # Convert withdrawals column to numeric
        df[withdrawals_col] = pd.to_numeric(df[withdrawals_col], errors='coerce')

        # Drop rows where withdrawals > 0
        df = df[df[withdrawals_col] <= 0]

    # Keywords to filter out (converted to lowercase)
    exclusion_keywords = [
        'interswitch transactions',
        'interswitch_pos terminal electronic money transfer levy',
        'nft/guaranty trust bank plc/bo/interswitchng/dr_goodsandservices',
        'nft/zenith bank plc/b/o/interswitchng/pr_goodsandservices',
        'stamp duty',
        'ltr dt',
        'import duty pymnt',
        'custom duty payment',
        'duty payment',
        'stampduty',
        'standing order',
        'unified payments',
        'sms alert charges',
        'sms fee',
        'sms notification charge',
        'goods inter transfer',
        'loans',
        'loan repayment',
        'usd payments',
        'fx sales',
        'settlement pay',
        'pp_ccril/settlement',
        'goods and services – interswitchng',
        'goodsandservices',
        'sms alert – sms',
        'usd', '$', 'mf2'
    ]
    exclusion_keywords = [word.lower() for word in exclusion_keywords]

    # Filter out transactions containing any of the exclusion keywords
    df_filtered = df[~df[narration_col].apply(lambda x: any(keyword in x for keyword in exclusion_keywords))]

    if amount_col:
        # Remove duplicate transactions (same amount, same narration)
        df_filtered = df_filtered[~df_filtered.duplicated(subset=[narration_col, amount_col], keep=False)]

    # **Creating the Pending Sheet**
    pending_keywords = ["cq", "cheque payments", "chq", "cheque deposits"]
    pending_keywords = [word.lower() for word in pending_keywords]

    # Select all transactions except excluded ones
    df_pending = df_filtered.copy()

    # Apply general pending keywords rule
    keyword_condition = df_pending[narration_col].apply(lambda x: any(kw in x for kw in pending_keywords))

    # Save both filtered and pending data to separate Excel files
    df_filtered.to_excel("filtered_transactions.xlsx", index=False)
    df_pending.to_excel("pending.xlsx", index=False)

    print("Filtered transactions saved to 'filtered_transactions.xlsx'")
    print("Pending transactions saved to 'pending_transactions.xlsx'")

    return df_filtered, df_pending

# Example usage
filter_transactions("Access (6).xlsx")

