// Decodes packets the server sent with the game's own codecs.
//
// Run by cmd/mcverify as a single-file source program, against the client jar:
//
//   java -cp <client jar and libraries> DecodePackets.java <packet dir> <results file>
//
// The packet directory holds one file per packet, named <phase>_<n>_<id>.bin, whose
// content is what the client's codec reads: the packet id as a VarInt, then the
// body.  <phase> is login, configuration or play.
//
// Results go to a file because bootstrapping the game redirects System.out into
// its own log format.  Each line is STATUS, a tab, the file name, a tab, a detail:
//
//   OK    decoded, and every byte was used
//   LEFT  decoded, but bytes were left over: the server sent more than the packet holds
//   FAIL  the codec threw: the packet is not the layout the client expects
//   SKIP  could not be checked here, for a reason given in the detail

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import java.io.File;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.util.Arrays;
import java.util.Map;
import net.minecraft.core.RegistryAccess;
import net.minecraft.network.ProtocolInfo;
import net.minecraft.network.RegistryFriendlyByteBuf;
import net.minecraft.network.protocol.configuration.ConfigurationProtocols;
import net.minecraft.network.protocol.game.GameProtocols;
import net.minecraft.network.protocol.login.LoginProtocols;

public class DecodePackets {
    public static void main(String[] args) throws Exception {
        net.minecraft.SharedConstants.tryDetectVersion();
        net.minecraft.server.Bootstrap.bootStrap();

        // Play packets that carry registry-backed values (the dimension in login, say)
        // need the registries the server sends during configuration, which this does not
        // build, so those are skipped when they ask for one.
        Map<String, ProtocolInfo<?>> protocols = Map.of(
            "login", LoginProtocols.CLIENTBOUND,
            "configuration", ConfigurationProtocols.CLIENTBOUND,
            "play", GameProtocols.CLIENTBOUND_TEMPLATE.bind(RegistryFriendlyByteBuf.decorator(RegistryAccess.EMPTY)));

        File[] files = new File(args[0]).listFiles();
        Arrays.sort(files);

        try (PrintWriter out = new PrintWriter(args[1])) {
            for (File file : files) {
                String phase = file.getName().substring(0, file.getName().indexOf('_'));
                ProtocolInfo<?> protocol = protocols.get(phase);
                if (protocol == null) {
                    out.println("SKIP\t" + file.getName() + "\tunknown phase " + phase);
                    continue;
                }

                byte[] bytes = Files.readAllBytes(file.toPath());
                ByteBuf buf = Unpooled.wrappedBuffer(bytes);
                try {
                    Object packet = protocol.codec().decode(buf);
                    String name = packet.getClass().getSimpleName();
                    if (buf.readableBytes() == 0) {
                        out.println("OK\t" + file.getName() + "\t" + name);
                    } else {
                        out.println("LEFT\t" + file.getName() + "\t" + name + ": " + buf.readableBytes()
                            + " of " + bytes.length + " bytes unread");
                    }
                } catch (Throwable t) {
                    Throwable root = t;
                    while (root.getCause() != null) {
                        root = root.getCause();
                    }
                    String detail = "at byte " + buf.readerIndex() + " of " + bytes.length + ": " + root;
                    String status = String.valueOf(root.getMessage()).startsWith("Missing registry") ? "SKIP" : "FAIL";
                    out.println(status + "\t" + file.getName() + "\t" + detail);
                }
            }
        }
    }
}
