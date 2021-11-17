using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Pomelo.DotNetClient;
using SimpleJson;
using System.Threading;

enum eState
{
    READY,
    QUERY_GATE,
    WORKING,
}

public class TestPomelo : MonoBehaviour
{
    string config_host = "127.0.0.1";//(www.xxx.com/127.0.0.1/::1/localhost etc.)
    int config_port = 30021;
    PomeloClient pclient;

    List<string> _messages = new List<string>();
    string _curInput;

    eState _state;

    // Start is called before the first frame update
    void Start()
    {
        Debug.Log("mainthread:" + Thread.CurrentThread.ManagedThreadId);
        createClient();
    }

    void prepare()
    {
        pclient.NetWorkStateChangedEvent += (state) =>
        {
            Debug.Log(state);
        };
    }

    void createClient()
    {
        pclient = new PomeloClient();
        pclient.msgPullMode = true;
        prepare();
    }

    void begin()
    {
        if(_state != eState.READY)
        {
            pclient.disconnect();
            createClient();
        }

        _state = eState.READY;
        startPomelo();
    }

    void startPomelo()
    {        
        _state = eState.QUERY_GATE;
        pclient.disconnect();
        connect(config_host, config_port);
    }

    void connect(string host, int port)
    {
        pclient.initClient(host, port, () =>
        {
            Debug.Log("init succ");
            //The user data is the handshake user params
            JsonObject user = new JsonObject();
            pclient.connect(user, data =>
            {
                Debug.Log("pclient.connect succ");
                onPomeloReady();
                //process handshake call back data
            });
        });
    }

    void onPomeloReady()
    {
        if(_state == eState.QUERY_GATE)
            queryGate();
        else
        {
            doLogin();
        }
    }

    void queryGate()
    {
        JsonObject args = new JsonObject();
        pclient.request("gate.gate.querygate", args, onQueryGateAck);
    }    

    void onQueryGateAck(JsonObject data)
    {
        Debug.Log(data);
        string ip = JsonDataExtension.AsStr(data, "ip", "127.0.0.1");
        string ports = JsonDataExtension.AsStr(data, "port", "");
        string[] subs = ports.Split(',');
        if (subs.Length < 2)
            return;
        int tcpPort = System.Convert.ToInt32(subs[1]);


        pclient.disconnect();
        createClient();
        
        _state = eState.WORKING;
        connect(ip, tcpPort);
        //pclient.NetWorkStateChangedEvent += (state) =>
        //{
        //    Debug.Log(state);
        //};
    }

    void doLogin()
    {
        JsonObject args = new JsonObject();
        args["name"] = "haha";
        pclient.request("gate.gate.login", args, onLoginAck);

        // 注册协议相关
        pclient.on("onMessage", onMessage);
    }

    void onLoginAck(JsonObject data)
    {
        Debug.Log("onLoginAck");
    }

    void onMessage(JsonObject data)
    {
        //Debug.Log(data);
        string name = JsonDataExtension.AsStr(data, "name", "");
        string content = JsonDataExtension.AsStr(data, "content", "");

        _messages.Add(name + ": " + content);

        while (_messages.Count > 10)
            _messages.RemoveAt(0);
    }

    // Update is called once per frame
    void Update()
    {
        pclient.PullMsgs();
    }

    private void OnGUI()
    {
        var oneHei = Screen.height / 32;
        if (GUILayout.Button("login", GUILayout.Height(oneHei)))
        {
            begin();
        }

        GUILayout.TextArea( makeMessages(), 
            GUILayout.Width(oneHei*20), 
            GUILayout.Height(oneHei * 10));

        GUILayout.BeginHorizontal();
        _curInput = GUILayout.TextField(_curInput, GUILayout.Width(oneHei * 10));
        if (GUILayout.Button("enter", GUILayout.Height(oneHei))) 
        {
            sendChats(_curInput);
            _curInput = "";
        }
        GUILayout.EndHorizontal();
    }

    private void sendChats(string content)
    {
        for (int i = 0; i < 100; i++)
            sendChat(content);
    }

    private void sendChat(string content)
    {
        if (_state == eState.WORKING)
        {
            JsonObject args = new JsonObject();
            args["name"] = "haha";
            args["content"] = content;
            pclient.request("chat.chat.sendchat", args, data => 
            {
                Debug.Log("send chat cb succ, thread:"+ Thread.CurrentThread.ManagedThreadId);
            });
        }
    }

    private string makeMessages()
    {
        string content = "";
        for(int i = 0; i < _messages.Count; i ++)
        {
            content += _messages[i] + "\n";            
        }
        return content;
    }
}
