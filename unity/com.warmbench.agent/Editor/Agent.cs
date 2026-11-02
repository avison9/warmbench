using UnityEditor;
using UnityEngine;

namespace Warmbench
{
    /// <summary>
    /// Entry point of a warmbench slot editor, started with
    /// <c>-executeMethod Warmbench.Agent.Serve</c>.
    /// </summary>
    public static class Agent
    {
        /// <summary>
        /// Will keep the editor alive and serve jobs from warmbenchd over a local socket.
        /// Not implemented yet: it logs an error and exits with code 2.
        /// </summary>
        public static void Serve()
        {
            Debug.LogError("warmbench: Agent.Serve is not implemented yet; see docs/POOL.md");
            EditorApplication.Exit(2);
        }
    }
}
