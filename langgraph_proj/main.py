import os
import openai
import requests
from groq import Groq
from typing import List
from dotenv import load_dotenv
from typing import TypedDict, Dict, Any
from groq.resources.chat import completions
from langgraph.graph import StateGraph, END, START
from langgraph.graph.message import add_messages


load_dotenv()
client = Groq(
    api_key=os.getenv("GROQ_API_KEY")
)


# State ki maa ka bhosda
class State(TypedDict):
    name: str
    email: str
    phone_number: str
    document_uploaded: bool
    document_uploaded_paths: List[str]
    document_processed: bool
    extracted_fields: List[Dict]
    payable_amount_calculated: float
    meeting_scheduled: bool


graph_builder = StateGraph(State)


def is_chitchat(user_input: str) -> bool:
    check_prompt = f"""Classify the following user message as either 'chitchat' or 'answer'.
                    If the message is a direct answer to a question (e.g., a name, email, or phone number), respond with 'answer'.
                    If the message is a casual or informal interaction (e.g., 'How are you?', 'Hello!', 'What's up?'), respond with 'chitchat'.
                    Message: "{user_input}"
                    If the user input is a question then it must go to chitchat.
                    Respond with only one word: chitchat or answer."""
    chat_completion = client.chat.completions.create(
        messages=[
            {"role": "system", "content": "You are a helpful assistant."},
            {"role": "user", "content": check_prompt}
        ],
        model="llama-3.3-70b-versatile",
    )
    content = chat_completion.choices[0].message.content
    if content is None:
        return False
    return "chitchat" in content


def ask_user(prompt: str, state: dict) -> str:
    system_prompt = f"""
        You are a helpful assistant that is collecting user information step by step.
        Here is what you currently know about the user:

        Name: {state.get('name') or 'Not provided'}
        Email: {state.get('email') or 'Not provided'}
        Phone Number: {state.get('phone_number') or 'Not provided'}

        Use this context to avoid repeating questions unnecessarily and to answer questions like "What was my name?".
        Only ask for missing fields or respond conversationally to general input.
    """
    chat_completion = client.chat.completions.create(
        messages=[
            {"role": "system", "content": system_prompt.strip()},
            {"role": "user", "content": prompt}
        ],
        model="llama-3.3-70b-versatile",
    )
    content = chat_completion.choices[0].message.content
    if content is None:
        return "No Response from LLM"
    return content.strip()


#Profile Collector node
def entity_collection(state: State):
    updated_state = dict(state)

    if not updated_state["name"]:
        while True:
            print("Assistant:", ask_user("Ask me for my name. (Just use that line to ask me that\
            and dont bother with unneccessary explanation)", updated_state))
            user_input = input("User: ").strip()

            if is_chitchat(user_input):
                print("Assistant:", ask_user(user_input, updated_state))
                continue

            if user_input:
                updated_state["name"] = user_input
                print("Assistant: Hello, " + user_input + ", nice to meet you!")
                break
            else:
                print(
                    "Assistant: I didn't catch your name. Could you please provide it?")

    if not updated_state["email"]:
        while True:
            print("Assistant:", ask_user(
                "Ask me for my email address. (Just use that line to ask me that)", updated_state))
            user_input = input("User: ").strip()

            if is_chitchat(user_input):
                print("Assistant:", ask_user(user_input, updated_state))
                continue

            if user_input:
                updated_state["email"] = user_input
                print("Assistant: Got it! Your email is: " + user_input)
                break
            else:
                print(
                    "Assistant: I didn't catch your email. Could you please provide it?")

    if not updated_state["phone_number"]:
        while True:
            print("Assistant:", ask_user(
                "Ask me for my phone number. (Just use that line to ask me that)", updated_state))
            user_input = input("User: ").strip()

            if is_chitchat(user_input):
                print("Assistant:", ask_user(user_input, updated_state))
                continue

            if user_input:
                updated_state["phone_number"] = user_input
                print("Assistant: Thanks! Your phone number is: " + user_input)
                break
            else:
                print(
                    "Assistant: I didn't catch your phone number. Could you please provide it?")

    return updated_state


#Document uploader node
def document_uploader(state: State):
    updated_state = dict(state)
    paths = []

    while True:
        print("Assistant:", ask_user("How many documents would you like to upload?", updated_state))
        user_input = input("User: ").strip()

        if is_chitchat(user_input):
            print("Assistant:", ask_user(user_input, updated_state))
            continue

        try:
            num_docs = int(user_input)
            if num_docs <= 0:
                raise ValueError
            break
        except ValueError:
            print("Assistant: Please enter a valid positive number.")

    for i in range(num_docs):
        while True:
            prompt = f"Please provide the full path to document {i + 1}:"
            print("Assistant:", ask_user(prompt, updated_state))
            doc_path = input("User: ").strip()

            if is_chitchat(doc_path):
                print("Assistant:", ask_user(doc_path, updated_state))
                continue

            if doc_path:
                paths.append(doc_path)
                break
            else:
                print("Assistant: I didn't catch that. Please enter the document path.")

    updated_state["document_uploaded"] = True
    updated_state["document_uploaded_paths"] = paths
    return updated_state

#Document extractor
# def extracted_fields_processor(state: State):
#     updated_state = dict(state)
#     extracted_results = []
#
#     api_endpoint = "http://localhost:8000/extract"
#
#     import pdb
#     pdb.set_trace()
#     for path in updated_state["document_uploaded_paths"]:
#         try:
#             with open(path, "rb") as f:
#                 files = {"file": f}
#                 response = requests.post(api_endpoint, files=files)
#                 if response.status_code == 200:
#                     extracted_results.append(response.json())
#                 else:
#                     extracted_results.append({"error": f"Failed to extract from {path}, status code {response.status_code}"})
#         except Exception as e:
#             extracted_results.append({"error": f"Exception for {path}: {str(e)}"})
#
#     updated_state["document_processed"] = True
#     updated_state["extracted_fields"] = extracted_results
#     return updated_state

def extracted_fields_processor(state: State):
    updated_state = dict(state)
    extracted_results = []

    api_endpoint = "https://mmc.lightinfosys.com/icrv2/inbound/4/?token=9879c248-38b9-4bea-86bd-f774f6a206ff"

    import pdb
    pdb.set_trace()
    for path in updated_state["document_uploaded_paths"]:
        try:
            with open(path, "rb") as f:
                files = {
                    'file_attachment': ('file', f, 'application/octet-stream')
                }
                payload = {
                    'model': 'templateless',
                    'detail': 'true'
                }
                response = requests.post(api_endpoint, data=payload, files=files)

                if response.status_code == 200:
                    extracted_results.append(response.json())
                else:
                    extracted_results.append({
                        "error": f"Failed to extract from {path}, status code {response.status_code}, response: {response.text}"
                    })
        except Exception as e:
            extracted_results.append({
                "error": f"Exception for {path}: {str(e)}"
            })

    updated_state["document_processed"] = True
    updated_state["extracted_fields"] = extracted_results
    return updated_state


def should_continue(state: State) -> str:
    if not state["name"] or not state["email"] or not state["phone_number"]:
        return "entity_collection"
    if not state["document_uploaded"]:
        return "document_uploader"
    return END


def should_continue_after_upload(state: State) -> str:
    if not state["document_processed"]:
        return "extracted_fields_processor"
    return END

graph_builder.add_node("entity_collection", entity_collection)
graph_builder.add_node("document_uploader", document_uploader)
graph_builder.add_node("extracted_fields_processor", extracted_fields_processor)
graph_builder.set_entry_point("entity_collection")

graph_builder.add_conditional_edges("entity_collection", should_continue)
graph_builder.add_conditional_edges("document_uploader", should_continue_after_upload)
# graph_builder.add_edge("document_uploader", END)
graph_builder.add_edge("extracted_fields_processor", END)


graph = graph_builder.compile()

initial_state: State = {
    "name": "",
    "email": "",
    "phone_number": "",
    "document_uploaded": False,
    "document_uploaded_paths": [],
    "document_processed": False,
    "extracted_fields": [],
    "payable_amount_calculated": 0.00,
    "meeting_scheduled": False,
}


def create_graph_loop(state: State):
    for event in graph.stream(state):
        for value in event.values():
            state = value
    print("\nUpdated State:", state)
    return state


final_state = create_graph_loop(initial_state)
