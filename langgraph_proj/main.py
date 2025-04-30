import os
import openai
from groq import Groq
from dotenv import load_dotenv
from langchain_ollama import OllamaLLM
from typing import TypedDict, Dict, Any
from groq.resources.chat import completions
from langgraph.graph import StateGraph, END, START
from langgraph.graph.message import add_messages



load_dotenv()
client = Groq(
    api_key=os.getenv("GROQ_API_KEY")
)



#State ki maa ka bhosda
class State(TypedDict):
    name: str
    email: str
    phone_number: str
    document_uploaded: bool
    document_processed: bool
    extracted_fields: Dict[str, Any]
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
            {"role": "user", "content":check_prompt}
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


def entity_collection(state: State):
    updated_state = dict(state)
    allowed_keys = ["name", "email", "phone number"]

    if not updated_state["name"]:
        while True:
            print("Assistant:", ask_user("Ask me for my name. (Just use that line to ask me that\
            and dont bother with unneccessary explanation)", updated_state))
            user_input = input("User: ").strip()

            if is_chitchat(user_input):
                print("Assistant:", ask_user(user_input,updated_state))
                continue

            if user_input:
                # if any(key in user_input for key in ["name", "email", "phone number"]):
                #     if "name" in user_input:
                #         print("Assistant: Your name is ", updated_state.get("name"))
                #     elif "email" in user_input:
                #         print("Assistant: Your email is ",updated_state.get("email"))
                #     elif "phone number" in user_input:
                #         print("Assistant: Your phone number is ", updated_state.get("phone_number"))
                updated_state["name"] = user_input
                print("Assistant: Hello, " + user_input + ", nice to meet you!")
                break
            else:
                print("Assistant: I didn't catch your name. Could you please provide it?")

    if not updated_state["email"]:
        while True:
            print("Assistant:", ask_user("Ask me for my email address. (Just use that line to ask me that)", updated_state))
            user_input = input("User: ").strip()

            if is_chitchat(user_input):
                print("Assistant:", ask_user(user_input, updated_state))
                continue

            if user_input:
                # if any(key in user_input for key in ["name", "email", "phone number"]):
                #     if "name" in user_input:
                #         print("Assistant: Your name is ", updated_state.get("name"))
                #     elif "email" in user_input:
                #         print("Assistant: Your email is ", updated_state.get("email"))
                #     elif "phone number" in user_input:
                #         print("Assistant: Your phone number is ", updated_state.get("phone_number"))
                # if any(word in user_input for word in user_input.split() if word not in allowed_keys):
                updated_state["email"] = user_input
                print("Assistant: Got it! Your email is: " + user_input)
                break
            else:
                print("Assistant: I didn't catch your email. Could you please provide it?")

    if not updated_state["phone_number"]:
        while True:
            print("Assistant:", ask_user("Ask me for my phone number. (Just use that line to ask me that)", updated_state))
            user_input = input("User: ").strip()

            if is_chitchat(user_input):
                print("Assistant:", ask_user(user_input, updated_state))
                continue

            if user_input:
                # if any(key in user_input for key in ["name", "email", "phone number"]):
                #     if "name" in user_input:
                #         print("Assistant: Your name is ", updated_state.get("name"))
                #     elif "email" in user_input:
                #         print("Assistant: Your email is ", updated_state.get("email"))
                #     elif "phone number" in user_input:
                #         print("Assistant: Your phone number is ", updated_state.get("phone_number"))
                # if any(word in user_input for word in user_input.split() if word not in allowed_keys):
                updated_state["phone_number"] = user_input
                print("Assistant: Thanks! Your phone number is: " + user_input)
                break
            else:
                print("Assistant: I didn't catch your phone number. Could you please provide it?")

    return updated_state



# def entity_collection(state: State):
#     updated_state = dict(state)
#
#     # Ask for name if not provided
#     if not updated_state["name"]:
#         while True:
#             print("Assistant:", ask_user("Ask me for my name. (Just use that line to ask me that)"))
#             user_input = input("User: ").strip()
#
#             if is_chitchat(user_input):
#                 print("Assistant:", ask_user(user_input))  # Let the assistant handle chit-chat
#                 continue
#
#             if user_input:  # Check if user provides a valid name
#                 updated_state["name"] = user_input
#                 print("Assistant: Hello, " + user_input + ", nice to meet you!")
#                 break  # Exit loop once name is collected
#             else:
#                 print("Assistant: I didn't catch your name. Could you please provide it?")
#
#     # Ask for email if not provided
#     if not updated_state["email"]:
#         while True:
#             print("Assistant:", ask_user("Ask me for my email address. (Just use that line to ask me that)"))
#             user_input = input("User: ").strip()
#
#             if is_chitchat(user_input):
#                 print("Assistant:", ask_user(user_input))  # Let the assistant handle chit-chat
#                 continue
#
#             if user_input:  # Check if user provides a valid email
#                 updated_state["email"] = user_input
#                 print("Assistant: Got it! Your email is: " + user_input)
#                 break  # Exit loop once email is collected
#             else:
#                 print("Assistant: I didn't catch your email. Could you please provide it?")
#
#     # Ask for phone number if not provided
#     if not updated_state["phone_number"]:
#         while True:
#             print("Assistant:", ask_user("Ask me for my phone number. (Just use that line to ask me that)"))
#             user_input = input("User: ").strip()
#
#             if is_chitchat(user_input):
#                 print("Assistant:", ask_user(user_input))  # Let the assistant handle chit-chat
#                 continue
#
#             if user_input:  # Check if user provides a valid phone number
#                 updated_state["phone_number"] = user_input
#                 print("Assistant: Thanks! Your phone number is: " + user_input)
#                 break  # Exit loop once phone number is collected
#             else:
#                 print("Assistant: I didn't catch your phone number. Could you please provide it?")
#
#     return updated_state
# #



graph_builder.add_node("entity_collection", entity_collection)
graph_builder.set_entry_point("entity_collection")


def should_continue(state: State) -> str:
    if not state["name"] or not state["email"] or not state["phone_number"]:
        return "entity_collection"
    return END

graph_builder.add_conditional_edges("entity_collection", should_continue)

graph = graph_builder.compile()

initial_state: State = {
    "name": "",
    "email": "",
    "phone_number": "",
    "document_uploaded": False,
    "document_processed": False,
    "extracted_fields": {},
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













